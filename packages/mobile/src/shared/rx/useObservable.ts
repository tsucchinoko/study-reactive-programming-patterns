import { useEffect, useRef, useState } from "react";
import { Observable, Subscription } from "rxjs";

/**
 * Observableの購読状態を表す判別共用体（Discriminated Union）。
 *
 * - `idle`    — まだsubscribeしていない初期状態
 * - `loading` — subscribe済みだがデータ未着
 * - `success` — データ到着済み
 * - `error`   — エラー発生
 */
type ObservableState<T> =
  | { status: "idle" }
  | { status: "loading" }
  | { status: "success"; data: T }
  | { status: "error"; error: unknown };

/**
 * RxJS ObservableをサブスクライブしてReactのステートにバインドする。
 * アンマウント時またはdepsの変更時に自動的にサブスクリプションを解除する。
 *
 * @param factory  Observableを返す関数（depsの変更時に再実行される）。
 * @param deps     依存配列（useEffectと同じセマンティクス）。
 */
export function useObservable<T>(
  factory: () => Observable<T>,
  deps: React.DependencyList,
): ObservableState<T> {
  const [state, setState] = useState<ObservableState<T>>({ status: "idle" });
  const subscriptionRef = useRef<Subscription | null>(null);

  useEffect(() => {
    setState({ status: "loading" });

    const observable$ = factory();
    subscriptionRef.current = observable$.subscribe({
      next: (data) => setState({ status: "success", data }),
      error: (error) => setState({ status: "error", error }),
    });

    return () => {
      subscriptionRef.current?.unsubscribe();
      subscriptionRef.current = null;
    };
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, deps);

  return state;
}
