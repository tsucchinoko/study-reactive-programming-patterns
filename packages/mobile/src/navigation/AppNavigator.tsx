import React from "react";
import { createNativeStackNavigator } from "@react-navigation/native-stack";
import { createBottomTabNavigator } from "@react-navigation/bottom-tabs";
import { OrderListScreen } from "../features/order/screens/OrderListScreen";
import { OrderDetailScreen } from "../features/order/screens/OrderDetailScreen";
import { RestaurantListScreen } from "../features/restaurant/screens/RestaurantListScreen";
import { Text } from "react-native";

export type RootStackParamList = {
  Tabs: undefined;
  OrderDetail: { orderId: string };
  // Keep flat for OrderListScreen navigation.navigate compatibility.
  OrderList: undefined;
  Restaurants: undefined;
};

type TabParamList = {
  OrderList: undefined;
  Restaurants: undefined;
};

const Stack = createNativeStackNavigator<RootStackParamList>();
const Tab = createBottomTabNavigator<TabParamList>();

function TabIcon({ label, focused }: { label: string; focused: boolean }) {
  return (
    <Text style={{ fontSize: 10, color: focused ? "#3B82F6" : "#9CA3AF" }}>
      {label}
    </Text>
  );
}

function TabNavigator() {
  return (
    <Tab.Navigator
      screenOptions={{
        headerStyle: { backgroundColor: "#fff" },
        headerTitleStyle: { fontWeight: "bold" },
        tabBarActiveTintColor: "#3B82F6",
        tabBarInactiveTintColor: "#9CA3AF",
      }}
    >
      <Tab.Screen
        name="OrderList"
        component={OrderListScreen}
        options={{
          title: "注文一覧",
          tabBarIcon: ({ focused }) => <TabIcon label="📋" focused={focused} />,
        }}
      />
      <Tab.Screen
        name="Restaurants"
        component={RestaurantListScreen}
        options={{
          title: "レストラン",
          tabBarIcon: ({ focused }) => <TabIcon label="🍽" focused={focused} />,
        }}
      />
    </Tab.Navigator>
  );
}

export function AppNavigator() {
  return (
    <Stack.Navigator>
      <Stack.Screen
        name="Tabs"
        component={TabNavigator}
        options={{ headerShown: false }}
      />
      <Stack.Screen
        name="OrderDetail"
        component={OrderDetailScreen}
        options={{ title: "注文詳細" }}
      />
    </Stack.Navigator>
  );
}
