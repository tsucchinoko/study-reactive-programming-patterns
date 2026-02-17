import { useMemo } from "react";
import { filter, map } from "rxjs";
import { apolloClient } from "../../../shared/graphql/client";
import {
  fromApolloSubscription,
  retryWithBackoff,
} from "../../../shared/rx/operators";
import { useObservable } from "../../../shared/rx/useObservable";
import { ORDER_STATUS_CHANGED } from "../graphql/operations";

type Order = {
  id: string;
  status: string;
  total: { display: string };
  items: Array<{
    name: string;
    quantity: number;
    subtotal: { display: string };
  }>;
  placedAt: string;
  confirmedAt: string | null;
  deliveredAt: string | null;
  cancelledAt: string | null;
  cancelReason: string | null;
};

/**
 * Subscribe to real-time order status changes via GraphQL Subscription + RxJS.
 *
 * The Apollo subscription is converted to an RxJS Observable so we can
 * apply operators like retryWithBackoff, and bridge it into React state
 * with useObservable.
 */
export function useOrderStream(orderId: string) {
  const order$ = useMemo(
    () =>
      fromApolloSubscription<{ data?: { orderStatusChanged?: Order } | null }>(
        apolloClient.subscribe({
          query: ORDER_STATUS_CHANGED,
          variables: { orderId },
        }),
      ).pipe(
        map((result) => result.data?.orderStatusChanged),
        filter((order): order is Order => order != null),
        retryWithBackoff(5, 1000),
      ),
    [orderId],
  );

  return useObservable(() => order$, [order$]);
}
