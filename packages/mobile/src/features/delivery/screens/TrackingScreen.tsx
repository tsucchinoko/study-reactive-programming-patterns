import { useQuery } from "@apollo/client/react";
import React from "react";
import {
  ActivityIndicator,
  ScrollView,
  StyleSheet,
  Text,
  View,
} from "react-native";
import type { NativeStackScreenProps } from "@react-navigation/native-stack";
import { OrderStatusBadge } from "../../order/components/OrderStatusBadge";
import { StatusTimeline } from "../../../shared/components/StatusTimeline";
import { useOrderStream } from "../../order/hooks/useOrderStream";
import { useDriverLocation } from "../hooks/useDriverLocation";
import { GET_ORDER } from "../../order/graphql/operations";
import { GET_DELIVERY_BY_ORDER } from "../graphql/operations";
import type { RootStackParamList } from "../../../navigation/AppNavigator";

type Props = NativeStackScreenProps<RootStackParamList, "Tracking">;

export function TrackingScreen({ route }: Props) {
  const { orderId } = route.params;

  // 注文データ (初回クエリ + リアルタイムストリーム)
  const { data: orderData, loading: orderLoading } = useQuery<{
    order: any;
  }>(GET_ORDER, { variables: { id: orderId } });
  const orderStream = useOrderStream(orderId);
  const order =
    orderStream.status === "success" ? orderStream.data : orderData?.order;

  // 配達アサイン情報
  const { data: deliveryData, loading: deliveryLoading } = useQuery<{
    deliveryByOrder: any;
  }>(GET_DELIVERY_BY_ORDER, {
    variables: { orderId },
    pollInterval: 5000,
  });
  const assignment = deliveryData?.deliveryByOrder;
  const driver = assignment?.driver;

  // ドライバー位置のリアルタイムストリーム (RxJS)
  const driverLocation = useDriverLocation(driver?.id ?? null);

  const loading = orderLoading && deliveryLoading && !order;

  if (loading) {
    return (
      <View style={styles.center}>
        <ActivityIndicator size="large" color="#3B82F6" />
      </View>
    );
  }

  return (
    <ScrollView style={styles.container}>
      {/* 注文ヘッダー */}
      {order && (
        <View style={styles.header}>
          <Text style={styles.orderId}>
            注文 #{order.id.slice(0, 8)}
          </Text>
          <OrderStatusBadge status={order.status} />
          {orderStream.status === "success" && (
            <View style={styles.liveBadge}>
              <View style={styles.liveDot} />
              <Text style={styles.liveText}>LIVE</Text>
            </View>
          )}
        </View>
      )}

      {/* ステータスタイムライン */}
      {order && (
        <View style={styles.section}>
          <Text style={styles.sectionTitle}>ステータス</Text>
          <StatusTimeline currentStatus={order.status} />
        </View>
      )}

      {/* ドライバー情報 */}
      {driver ? (
        <View style={styles.section}>
          <Text style={styles.sectionTitle}>配達員</Text>
          <View style={styles.driverCard}>
            <View style={styles.driverAvatar}>
              <Text style={styles.driverAvatarText}>
                {driver.name.charAt(0)}
              </Text>
            </View>
            <View style={styles.driverInfo}>
              <Text style={styles.driverName}>{driver.name}</Text>
              <Text style={styles.driverPhone}>{driver.phone}</Text>
              <Text style={styles.driverStatus}>{driver.status}</Text>
            </View>
          </View>
        </View>
      ) : (
        <View style={styles.section}>
          <Text style={styles.sectionTitle}>配達員</Text>
          <Text style={styles.waitingText}>
            配達員のアサインを待っています...
          </Text>
        </View>
      )}

      {/* ドライバー位置 (リアルタイム) */}
      <View style={styles.section}>
        <Text style={styles.sectionTitle}>配達員の位置</Text>
        {driverLocation.status === "success" && driverLocation.data ? (
          <View style={styles.locationCard}>
            <View style={styles.locationLiveIndicator}>
              <View style={styles.locationLiveDot} />
              <Text style={styles.locationLiveText}>リアルタイム追跡中</Text>
            </View>
            <View style={styles.coordRow}>
              <Text style={styles.coordLabel}>緯度</Text>
              <Text style={styles.coordValue}>
                {driverLocation.data.latitude.toFixed(6)}
              </Text>
            </View>
            <View style={styles.coordRow}>
              <Text style={styles.coordLabel}>経度</Text>
              <Text style={styles.coordValue}>
                {driverLocation.data.longitude.toFixed(6)}
              </Text>
            </View>
            <Text style={styles.timestampText}>
              最終更新:{" "}
              {new Date(driverLocation.data.timestamp).toLocaleTimeString(
                "ja-JP",
              )}
            </Text>
          </View>
        ) : driverLocation.status === "loading" ? (
          <View style={styles.locationWaiting}>
            <ActivityIndicator size="small" color="#8B5CF6" />
            <Text style={styles.waitingText}>位置情報を受信中...</Text>
          </View>
        ) : (
          <Text style={styles.waitingText}>
            位置情報はまだ利用できません
          </Text>
        )}
      </View>

      {/* 配達ステータス */}
      {assignment && (
        <View style={styles.section}>
          <Text style={styles.sectionTitle}>配達状況</Text>
          <View style={styles.deliveryTimeline}>
            <TimelineItem
              label="アサイン"
              time={assignment.assignedAt}
              active={true}
            />
            <TimelineItem
              label="ピックアップ"
              time={assignment.pickedUpAt}
              active={!!assignment.pickedUpAt}
            />
            <TimelineItem
              label="配達完了"
              time={assignment.deliveredAt}
              active={!!assignment.deliveredAt}
            />
          </View>
        </View>
      )}
    </ScrollView>
  );
}

function TimelineItem({
  label,
  time,
  active,
}: {
  label: string;
  time: string | null;
  active: boolean;
}) {
  return (
    <View style={styles.timelineItem}>
      <View
        style={[
          styles.timelineDot,
          active ? styles.timelineDotActive : styles.timelineDotInactive,
        ]}
      />
      <View style={styles.timelineContent}>
        <Text
          style={[
            styles.timelineLabel,
            active ? styles.timelineLabelActive : styles.timelineLabelInactive,
          ]}
        >
          {label}
        </Text>
        {time && (
          <Text style={styles.timelineTime}>
            {new Date(time).toLocaleTimeString("ja-JP")}
          </Text>
        )}
      </View>
    </View>
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
  // ドライバーカード
  driverCard: {
    flexDirection: "row",
    alignItems: "center",
    gap: 16,
  },
  driverAvatar: {
    width: 48,
    height: 48,
    borderRadius: 24,
    backgroundColor: "#8B5CF6",
    justifyContent: "center",
    alignItems: "center",
  },
  driverAvatarText: {
    fontSize: 20,
    fontWeight: "bold",
    color: "#fff",
  },
  driverInfo: {
    flex: 1,
  },
  driverName: {
    fontSize: 16,
    fontWeight: "bold",
    color: "#1F2937",
  },
  driverPhone: {
    fontSize: 13,
    color: "#6B7280",
    marginTop: 2,
  },
  driverStatus: {
    fontSize: 12,
    color: "#8B5CF6",
    fontWeight: "600",
    marginTop: 2,
  },
  waitingText: {
    fontSize: 14,
    color: "#9CA3AF",
    marginLeft: 8,
  },
  // 位置カード
  locationCard: {
    backgroundColor: "#F5F3FF",
    borderRadius: 12,
    padding: 16,
    gap: 8,
  },
  locationLiveIndicator: {
    flexDirection: "row",
    alignItems: "center",
    marginBottom: 4,
  },
  locationLiveDot: {
    width: 8,
    height: 8,
    borderRadius: 4,
    backgroundColor: "#8B5CF6",
    marginRight: 6,
  },
  locationLiveText: {
    fontSize: 12,
    fontWeight: "bold",
    color: "#8B5CF6",
  },
  coordRow: {
    flexDirection: "row",
    justifyContent: "space-between",
    alignItems: "center",
  },
  coordLabel: {
    fontSize: 13,
    color: "#6B7280",
  },
  coordValue: {
    fontSize: 15,
    fontWeight: "600",
    color: "#1F2937",
    fontVariant: ["tabular-nums"],
  },
  timestampText: {
    fontSize: 11,
    color: "#9CA3AF",
    marginTop: 4,
  },
  locationWaiting: {
    flexDirection: "row",
    alignItems: "center",
    gap: 8,
  },
  // 配達タイムライン
  deliveryTimeline: {
    gap: 16,
  },
  timelineItem: {
    flexDirection: "row",
    alignItems: "center",
    gap: 12,
  },
  timelineDot: {
    width: 12,
    height: 12,
    borderRadius: 6,
  },
  timelineDotActive: {
    backgroundColor: "#8B5CF6",
  },
  timelineDotInactive: {
    backgroundColor: "#D1D5DB",
  },
  timelineContent: {
    flex: 1,
    flexDirection: "row",
    justifyContent: "space-between",
    alignItems: "center",
  },
  timelineLabel: {
    fontSize: 14,
  },
  timelineLabelActive: {
    color: "#1F2937",
    fontWeight: "600",
  },
  timelineLabelInactive: {
    color: "#9CA3AF",
  },
  timelineTime: {
    fontSize: 12,
    color: "#6B7280",
  },
});
