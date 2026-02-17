import React from "react";
import { StyleSheet, Text, View } from "react-native";

const STATUSES = [
  "CREATED",
  "CONFIRMED",
  "PREPARING",
  "READY",
  "PICKED_UP",
  "DELIVERING",
  "DELIVERED",
] as const;

const STATUS_LABELS: Record<string, string> = {
  CREATED: "注文受付",
  CONFIRMED: "確認済み",
  PREPARING: "調理中",
  READY: "準備完了",
  PICKED_UP: "ピックアップ",
  DELIVERING: "配達中",
  DELIVERED: "配達完了",
  CANCELLED: "キャンセル",
};

type Props = {
  currentStatus: string;
};

export function StatusTimeline({ currentStatus }: Props) {
  if (currentStatus === "CANCELLED") {
    return (
      <View style={styles.container}>
        <View style={styles.step}>
          <View style={[styles.dot, styles.dotCancelled]} />
          <Text style={[styles.label, styles.labelCancelled]}>
            キャンセル済み
          </Text>
        </View>
      </View>
    );
  }

  const currentIdx = STATUSES.indexOf(currentStatus as (typeof STATUSES)[number]);

  return (
    <View style={styles.container}>
      {STATUSES.map((status, idx) => {
        const isPast = idx <= currentIdx;
        const isCurrent = idx === currentIdx;

        return (
          <View key={status} style={styles.step}>
            <View style={styles.dotRow}>
              <View
                style={[
                  styles.dot,
                  isPast && styles.dotActive,
                  isCurrent && styles.dotCurrent,
                ]}
              />
              {idx < STATUSES.length - 1 && (
                <View
                  style={[styles.line, isPast && styles.lineActive]}
                />
              )}
            </View>
            <Text
              style={[
                styles.label,
                isPast && styles.labelActive,
                isCurrent && styles.labelCurrent,
              ]}
            >
              {STATUS_LABELS[status]}
            </Text>
          </View>
        );
      })}
    </View>
  );
}

const styles = StyleSheet.create({
  container: {
    paddingVertical: 16,
    paddingHorizontal: 8,
  },
  step: {
    flexDirection: "row",
    alignItems: "center",
    minHeight: 40,
  },
  dotRow: {
    alignItems: "center",
    width: 24,
  },
  dot: {
    width: 12,
    height: 12,
    borderRadius: 6,
    backgroundColor: "#D1D5DB",
  },
  dotActive: {
    backgroundColor: "#3B82F6",
  },
  dotCurrent: {
    backgroundColor: "#3B82F6",
    width: 16,
    height: 16,
    borderRadius: 8,
    borderWidth: 3,
    borderColor: "#BFDBFE",
  },
  dotCancelled: {
    backgroundColor: "#EF4444",
    width: 16,
    height: 16,
    borderRadius: 8,
  },
  line: {
    width: 2,
    height: 20,
    backgroundColor: "#D1D5DB",
  },
  lineActive: {
    backgroundColor: "#3B82F6",
  },
  label: {
    marginLeft: 12,
    fontSize: 14,
    color: "#9CA3AF",
  },
  labelActive: {
    color: "#374151",
  },
  labelCurrent: {
    fontWeight: "700",
    color: "#3B82F6",
  },
  labelCancelled: {
    fontWeight: "700",
    color: "#EF4444",
    marginLeft: 12,
    fontSize: 14,
  },
});
