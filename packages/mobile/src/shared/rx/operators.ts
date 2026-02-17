import { Observable, retry, timer } from "rxjs";

/**
 * Retry with exponential backoff.
 *
 * Example: retryWithBackoff(3, 1000)
 *   → 1st retry after ~1s, 2nd after ~2s, 3rd after ~4s
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

/** Minimal subset of the Observable-like interface Apollo Client returns. */
type ObservableLike<T> = {
  subscribe(observer: {
    next?: (value: T) => void;
    error?: (error: any) => void;
    complete?: () => void;
  }): { unsubscribe: () => void };
};

/**
 * Create an RxJS Observable from an Apollo subscription observable.
 * Apollo's Observable is not a true RxJS Observable — this bridges the gap.
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
