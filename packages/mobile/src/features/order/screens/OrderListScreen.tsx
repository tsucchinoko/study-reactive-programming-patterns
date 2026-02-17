import { useQuery } from "@apollo/client/react";
import React, { useCallback } from "react";
import {
  ActivityIndicator,
  FlatList,
  StyleSheet,
  Text,
  View,
} from "react-native";
import type { NativeStackScreenProps } from "@react-navigation/native-stack";
import { OrderCard } from "../components/OrderCard";
import { GET_ORDERS } from "../graphql/operations";
import type { RootStackParamList } from "../../../navigation/AppNavigator";

type Props = NativeStackScreenProps<RootStackParamList, "OrderList">;

export function OrderListScreen({ navigation }: Props) {
  const { data, loading, error, refetch } = useQuery<{ orders: any[] }>(
    GET_ORDERS,
    {
      variables: { limit: 50 },
      pollInterval: 5000,
    },
  );

  const handlePress = useCallback(
    (orderId: string) => {
      navigation.navigate("OrderDetail", { orderId });
    },
    [navigation],
  );

  if (loading && !data) {
    return (
      <View style={styles.center}>
        <ActivityIndicator size="large" color="#3B82F6" />
      </View>
    );
  }

  if (error) {
    return (
      <View style={styles.center}>
        <Text style={styles.error}>エラー: {error.message}</Text>
      </View>
    );
  }

  const orders = data?.orders ?? [];

  return (
    <FlatList
      data={orders}
      keyExtractor={(item) => item.id}
      renderItem={({ item }) => (
        <OrderCard order={item} onPress={handlePress} />
      )}
      contentContainerStyle={styles.list}
      onRefresh={() => refetch()}
      refreshing={loading}
      ListEmptyComponent={
        <View style={styles.center}>
          <Text style={styles.empty}>注文がありません</Text>
          <Text style={styles.emptyHint}>
            Simulatorを起動して注文を生成してください
          </Text>
        </View>
      }
    />
  );
}

const styles = StyleSheet.create({
  list: {
    paddingVertical: 8,
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
  empty: {
    fontSize: 16,
    fontWeight: "bold",
    color: "#6B7280",
  },
  emptyHint: {
    fontSize: 13,
    color: "#9CA3AF",
    marginTop: 4,
  },
});
