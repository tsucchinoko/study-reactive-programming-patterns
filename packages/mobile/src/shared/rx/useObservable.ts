import { useEffect, useRef, useState } from "react";
import { Observable, Subscription } from "rxjs";

type ObservableState<T> =
  | { status: "idle" }
  | { status: "loading" }
  | { status: "success"; data: T }
  | { status: "error"; error: unknown };

/**
 * Subscribe to an RxJS Observable and bind its emissions to React state.
 * Automatically unsubscribes on unmount or when deps change.
 *
 * @param factory  Function returning the Observable (called when deps change).
 * @param deps     Dependency array (same semantics as useEffect).
 */
export function useObservable<T>(
  factory: () => Observable<T>,
  deps: React.DependencyList
): ObservableState<T> {
  const [state, setState] = useState<ObservableState<T>>({ status: "idle" });
  const subscriptionRef = useRef<Subscription | null>(null);

  useEffect(() => {
    setState({ status: "loading" });

    const observable = factory();
    subscriptionRef.current = observable.subscribe({
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
