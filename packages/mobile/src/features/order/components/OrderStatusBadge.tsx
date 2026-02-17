import React from "react";
import { StyleSheet, Text, View } from "react-native";

const STATUS_CONFIG: Record<string, { label: string; color: string }> = {
  CREATED: { label: "注文受付", color: "#6B7280" },
  CONFIRMED: { label: "確認済み", color: "#3B82F6" },
  PREPARING: { label: "調理中", color: "#F59E0B" },
  READY: { label: "準備完了", color: "#10B981" },
  PICKED_UP: { label: "ピックアップ済", color: "#8B5CF6" },
  DELIVERING: { label: "配達中", color: "#EC4899" },
  DELIVERED: { label: "配達完了", color: "#059669" },
  CANCELLED: { label: "キャンセル", color: "#EF4444" },
};

type Props = {
  status: string;
};

export function OrderStatusBadge({ status }: Props) {
  const config = STATUS_CONFIG[status] ?? {
    label: status,
    color: "#6B7280",
  };

  return (
    <View style={[styles.badge, { backgroundColor: config.color + "20" }]}>
      <View style={[styles.dot, { backgroundColor: config.color }]} />
      <Text style={[styles.label, { color: config.color }]}>
        {config.label}
      </Text>
    </View>
  );
}

const styles = StyleSheet.create({
  badge: {
    flexDirection: "row",
    alignItems: "center",
    paddingHorizontal: 10,
    paddingVertical: 4,
    borderRadius: 12,
    alignSelf: "flex-start",
  },
  dot: {
    width: 6,
    height: 6,
    borderRadius: 3,
    marginRight: 6,
  },
  label: {
    fontSize: 12,
    fontWeight: "bold",
  },
});
