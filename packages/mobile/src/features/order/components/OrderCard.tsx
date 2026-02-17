import React from "react";
import { StyleSheet, Text, TouchableOpacity, View } from "react-native";
import { OrderStatusBadge } from "./OrderStatusBadge";

type OrderSummary = {
  id: string;
  status: string;
  total: { display: string };
  items: Array<{ name: string; quantity: number }>;
  placedAt: string;
};

type Props = {
  order: OrderSummary;
  onPress: (orderId: string) => void;
};

export function OrderCard({ order, onPress }: Props) {
  const itemSummary = order.items
    .slice(0, 3)
    .map((i) => `${i.name} x${i.quantity}`)
    .join("、");
  const more = order.items.length > 3 ? ` 他${order.items.length - 3}品` : "";
  const placedAt = new Date(order.placedAt).toLocaleString("ja-JP", {
    hour: "2-digit",
    minute: "2-digit",
  });

  return (
    <TouchableOpacity
      style={styles.card}
      onPress={() => onPress(order.id)}
      activeOpacity={0.7}
    >
      <View style={styles.header}>
        <Text style={styles.orderId}>#{order.id.slice(0, 8)}</Text>
        <OrderStatusBadge status={order.status} />
      </View>
      <Text style={styles.items} numberOfLines={1}>
        {itemSummary}
        {more}
      </Text>
      <View style={styles.footer}>
        <Text style={styles.time}>{placedAt}</Text>
        <Text style={styles.total}>{order.total.display}</Text>
      </View>
    </TouchableOpacity>
  );
}

const styles = StyleSheet.create({
  card: {
    backgroundColor: "#fff",
    borderRadius: 12,
    padding: 16,
    marginHorizontal: 16,
    marginVertical: 6,
    shadowColor: "#000",
    shadowOffset: { width: 0, height: 1 },
    shadowOpacity: 0.08,
    shadowRadius: 4,
    elevation: 2,
  },
  header: {
    flexDirection: "row",
    justifyContent: "space-between",
    alignItems: "center",
    marginBottom: 8,
  },
  orderId: {
    fontSize: 14,
    fontWeight: "700",
    color: "#1F2937",
  },
  items: {
    fontSize: 13,
    color: "#6B7280",
    marginBottom: 8,
  },
  footer: {
    flexDirection: "row",
    justifyContent: "space-between",
    alignItems: "center",
  },
  time: {
    fontSize: 12,
    color: "#9CA3AF",
  },
  total: {
    fontSize: 15,
    fontWeight: "700",
    color: "#1F2937",
  },
});
