import { Observable, retry, timer } from "rxjs";

/**
 * 指数バックオフ付きリトライ。
 *
 * 例: retryWithBackoff(3, 1000)
 *   → 1回目: 約1秒後、2回目: 約2秒後、3回目: 約4秒後
 */
export const retryWithBackoff = <T>(
  maxRetries: number,
  initialDelayMs: number,
) =>
  retry<T>({
    count: maxRetries,
    delay: (_error, retryCount) =>
      timer(initialDelayMs * Math.pow(2, retryCount - 1)),
  });

/** Apollo Clientが返すObservableライクなインターフェースの最小サブセット。 */
type ObservableLike<T> = {
  subscribe(observer: {
    next?: (value: T) => void;
    error?: (error: any) => void;
    complete?: () => void;
  }): { unsubscribe: () => void };
};

/**
 * ApolloのSubscription ObservableからRxJS Observableを生成する。
 * ApolloのObservableは真のRxJS Observableではないため、このブリッジが必要。
 */
export const fromApolloSubscription = <T>(
  apolloObservable: ObservableLike<T>,
): Observable<T> =>
  new Observable<T>((subscriber) => {
    const subscription = apolloObservable.subscribe({
      next: (value: T) => subscriber.next(value),
      error: (err: any) => subscriber.error(err),
      complete: () => subscriber.complete(),
    });
    return () => subscription.unsubscribe();
  });
