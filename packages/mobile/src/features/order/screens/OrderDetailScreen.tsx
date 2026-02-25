import { useQuery } from "@apollo/client/react";
import React from "react";
import {
  ActivityIndicator,
  Pressable,
  ScrollView,
  StyleSheet,
  Text,
  View,
} from "react-native";
import type { NativeStackScreenProps } from "@react-navigation/native-stack";
import { useNavigation } from "@react-navigation/native";
import type { NativeStackNavigationProp } from "@react-navigation/native-stack";
import { OrderStatusBadge } from "../components/OrderStatusBadge";
import { StatusTimeline } from "../../../shared/components/StatusTimeline";
import { useOrderStream } from "../hooks/useOrderStream";
import { GET_ORDER } from "../graphql/operations";
import type { RootStackParamList } from "../../../navigation/AppNavigator";

type Props = NativeStackScreenProps<RootStackParamList, "OrderDetail">;

const TRACKABLE_STATUSES = ["READY", "PICKED_UP", "DELIVERING"];

export function OrderDetailScreen({ route }: Props) {
  const { orderId } = route.params;
  const navigation =
    useNavigation<NativeStackNavigationProp<RootStackParamList>>();

  // Initial data via query.
  const { data, loading, error } = useQuery<{ order: any }>(GET_ORDER, {
    variables: { id: orderId },
  });

  // Real-time updates via subscription → RxJS.
  const stream = useOrderStream(orderId);

  // Prefer the latest streamed data, fallback to query data.
  const order = stream.status === "success" ? stream.data : data?.order;

  if (loading && !order) {
    return (
      <View style={styles.center}>
        <ActivityIndicator size="large" color="#3B82F6" />
      </View>
    );
  }

  if (error && !order) {
    return (
      <View style={styles.center}>
        <Text style={styles.error}>エラー: {error.message}</Text>
      </View>
    );
  }

  if (!order) {
    return (
      <View style={styles.center}>
        <Text style={styles.error}>注文が見つかりません</Text>
      </View>
    );
  }

  return (
    <ScrollView style={styles.container}>
      {/* Header */}
      <View style={styles.header}>
        <Text style={styles.orderId}>注文 #{order.id.slice(0, 8)}</Text>
        <OrderStatusBadge status={order.status} />
        {stream.status === "success" && (
          <View style={styles.liveBadge}>
            <View style={styles.liveDot} />
            <Text style={styles.liveText}>LIVE</Text>
          </View>
        )}
      </View>

      {/* Status Timeline */}
      <View style={styles.section}>
        <Text style={styles.sectionTitle}>ステータス</Text>
        <StatusTimeline currentStatus={order.status} />
      </View>

      {/* Order Items */}
      <View style={styles.section}>
        <Text style={styles.sectionTitle}>注文内容</Text>
        {order.items.map((item: any, idx: number) => (
          <View key={idx} style={styles.itemRow}>
            <View style={styles.itemInfo}>
              <Text style={styles.itemName}>{item.name}</Text>
              <Text style={styles.itemQty}>x{item.quantity}</Text>
            </View>
            <Text style={styles.itemPrice}>{item.subtotal.display}</Text>
          </View>
        ))}
        <View style={styles.totalRow}>
          <Text style={styles.totalLabel}>合計</Text>
          <Text style={styles.totalValue}>{order.total.display}</Text>
        </View>
      </View>

      {/* Cancel Reason */}
      {order.cancelReason && (
        <View style={styles.section}>
          <Text style={styles.sectionTitle}>キャンセル理由</Text>
          <Text style={styles.cancelReason}>{order.cancelReason}</Text>
        </View>
      )}

      {/* Tracking Button */}
      {TRACKABLE_STATUSES.includes(order.status) && (
        <View style={styles.section}>
          <Pressable
            style={styles.trackingButton}
            onPress={() => navigation.navigate("Tracking", { orderId })}
          >
            <Text style={styles.trackingButtonText}>配達追跡</Text>
          </Pressable>
        </View>
      )}
    </ScrollView>
  );
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: "#F9FAFB",
  },
  center: {
    flex: 1,
    justifyContent: "center",
    alignItems: "center",
    padding: 32,
  },
  error: {
    color: "#EF4444",
    fontSize: 14,
  },
  header: {
    backgroundColor: "#fff",
    padding: 20,
    flexDirection: "row",
    alignItems: "center",
    gap: 12,
    borderBottomWidth: 1,
    borderBottomColor: "#F3F4F6",
  },
  orderId: {
    fontSize: 18,
    fontWeight: "bold",
    color: "#1F2937",
    flex: 1,
  },
  liveBadge: {
    flexDirection: "row",
    alignItems: "center",
    backgroundColor: "#FEE2E2",
    paddingHorizontal: 8,
    paddingVertical: 2,
    borderRadius: 8,
  },
  liveDot: {
    width: 6,
    height: 6,
    borderRadius: 3,
    backgroundColor: "#EF4444",
    marginRight: 4,
  },
  liveText: {
    fontSize: 10,
    fontWeight: "bold",
    color: "#EF4444",
  },
  section: {
    backgroundColor: "#fff",
    marginTop: 8,
    padding: 20,
  },
  sectionTitle: {
    fontSize: 14,
    fontWeight: "bold",
    color: "#374151",
    marginBottom: 12,
  },
  itemRow: {
    flexDirection: "row",
    justifyContent: "space-between",
    alignItems: "center",
    paddingVertical: 8,
    borderBottomWidth: 1,
    borderBottomColor: "#F3F4F6",
  },
  itemInfo: {
    flexDirection: "row",
    alignItems: "center",
    gap: 8,
    flex: 1,
  },
  itemName: {
    fontSize: 14,
    color: "#374151",
  },
  itemQty: {
    fontSize: 12,
    color: "#9CA3AF",
  },
  itemPrice: {
    fontSize: 14,
    fontWeight: "bold",
    color: "#1F2937",
  },
  totalRow: {
    flexDirection: "row",
    justifyContent: "space-between",
    alignItems: "center",
    paddingTop: 12,
    marginTop: 4,
  },
  totalLabel: {
    fontSize: 15,
    fontWeight: "bold",
    color: "#374151",
  },
  totalValue: {
    fontSize: 18,
    fontWeight: "bold",
    color: "#1F2937",
  },
  cancelReason: {
    fontSize: 14,
    color: "#EF4444",
  },
  trackingButton: {
    backgroundColor: "#8B5CF6",
    paddingVertical: 14,
    borderRadius: 12,
    alignItems: "center",
  },
  trackingButtonText: {
    color: "#fff",
    fontSize: 16,
    fontWeight: "bold",
  },
});
