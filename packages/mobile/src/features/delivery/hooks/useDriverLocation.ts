import { useMemo } from "react";
import { distinctUntilChanged, filter, map, scan, throttleTime } from "rxjs";
import { apolloClient } from "../../../shared/graphql/client";
import {
  fromApolloSubscription,
  retryWithBackoff,
} from "../../../shared/rx/operators";
import { useObservable } from "../../../shared/rx/useObservable";
import { DRIVER_LOCATION_UPDATED } from "../graphql/operations";

export type DriverLocation = {
  driverId: string;
  latitude: number;
  longitude: number;
  timestamp: string;
};

/**
 * ドライバー位置をリアルタイムにストリーミングする。
 *
 * RxJSパイプラインで以下のリアクティブパターンを適用:
 *   1. throttleTime(1000) — 高頻度な位置更新を1秒に1回に間引き
 *   2. scan (lerp)        — 前回位置と現在位置の線形補間でスムーズな移動を表現
 *   3. distinctUntilChanged — 同一位置の重複イベントを排除
 *   4. retryWithBackoff    — 接続切断時の指数バックオフ再接続
 */
export function useDriverLocation(driverId: string | null) {
  const location$ = useMemo(() => {
    if (!driverId) return null;

    return fromApolloSubscription<{
      data?: { driverLocationUpdated?: DriverLocation } | null;
    }>(
      apolloClient.subscribe({
        query: DRIVER_LOCATION_UPDATED,
        variables: { driverId },
      }),
    ).pipe(
      // GraphQL レスポンスから位置データを抽出
      map((result) => result.data?.driverLocationUpdated),
      filter((loc): loc is DriverLocation => loc != null),

      // 1. throttleTime: 高頻度な位置更新を1秒に1回に間引く
      //    サーバーは1秒ごとに位置を送信するが、ネットワーク遅延でバーストする可能性がある
      throttleTime(1000),

      // 2. scan + lerp: 前回位置と現在位置の線形補間
      //    急激な位置ジャンプを抑え、スムーズな移動を表現する
      scan(
        (prev, current) => ({
          ...current,
          latitude: lerp(prev.latitude, current.latitude, 0.7),
          longitude: lerp(prev.longitude, current.longitude, 0.7),
        }),
      ),

      // 3. distinctUntilChanged: 緯度・経度が変化していない場合は再レンダリングしない
      distinctUntilChanged(
        (a, b) => a.latitude === b.latitude && a.longitude === b.longitude,
      ),

      // 4. retryWithBackoff: 接続切断時に指数バックオフで再接続
      retryWithBackoff(5, 1000),
    );
  }, [driverId]);

  return useObservable(() => location$, [location$]);
}

/** 線形補間 (Linear Interpolation) */
function lerp(a: number, b: number, t: number): number {
  return a + (b - a) * t;
}
