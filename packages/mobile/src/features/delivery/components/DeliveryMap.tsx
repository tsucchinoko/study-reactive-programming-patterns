import React, { useRef, useEffect } from "react";
import { StyleSheet, View, Text } from "react-native";
import { WebView } from "react-native-webview";
import type { DriverLocation } from "../hooks/useDriverLocation";

type Props = {
  driverLocation: DriverLocation;
  driverName?: string;
};

/**
 * Leaflet (OpenStreetMap) を WebView で描画し、配達員の位置をリアルタイムに表示する。
 *
 * react-native-maps はネイティブビルドが必要なため、
 * Expo Go でも動作する WebView + Leaflet を採用。
 */
export function DeliveryMap({ driverLocation, driverName }: Props) {
  const webViewRef = useRef<WebView>(null);
  const initializedRef = useRef(false);

  // 位置が更新されるたびに WebView 側の JavaScript を呼び出してマーカーを移動
  useEffect(() => {
    if (!initializedRef.current) return;
    webViewRef.current?.injectJavaScript(`
      updateDriverLocation(${driverLocation.latitude}, ${driverLocation.longitude}, "${new Date(driverLocation.timestamp).toLocaleTimeString("ja-JP")}");
      true;
    `);
  }, [driverLocation.latitude, driverLocation.longitude, driverLocation.timestamp]);

  const html = buildMapHtml(driverLocation, driverName);

  return (
    <View style={styles.container}>
      <View style={styles.liveIndicator}>
        <View style={styles.liveDot} />
        <Text style={styles.liveText}>リアルタイム追跡中</Text>
      </View>
      <View style={styles.mapContainer}>
        <WebView
          ref={webViewRef}
          style={styles.map}
          source={{ html }}
          scrollEnabled={false}
          onLoad={() => {
            initializedRef.current = true;
          }}
        />
      </View>
    </View>
  );
}

function buildMapHtml(location: DriverLocation, driverName?: string): string {
  const name = driverName ?? "配達員";
  return `
<!DOCTYPE html>
<html>
<head>
  <meta name="viewport" content="width=device-width, initial-scale=1.0, maximum-scale=1.0, user-scalable=no">
  <link rel="stylesheet" href="https://unpkg.com/leaflet@1.9.4/dist/leaflet.css" />
  <script src="https://unpkg.com/leaflet@1.9.4/dist/leaflet.js"></script>
  <style>
    * { margin: 0; padding: 0; }
    #map { width: 100%; height: 100vh; }
  </style>
</head>
<body>
  <div id="map"></div>
  <script>
    var map = L.map('map', { zoomControl: false }).setView([${location.latitude}, ${location.longitude}], 16);
    L.tileLayer('https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png', {
      attribution: '&copy; OpenStreetMap contributors'
    }).addTo(map);

    var marker = L.marker([${location.latitude}, ${location.longitude}]).addTo(map);
    marker.bindPopup('<b>${name}</b>').openPopup();

    function updateDriverLocation(lat, lng, time) {
      marker.setLatLng([lat, lng]);
      marker.setPopupContent('<b>${name}</b><br>最終更新: ' + time);
      map.panTo([lat, lng]);
    }
  </script>
</body>
</html>
`;
}

const styles = StyleSheet.create({
  container: {
    gap: 8,
  },
  liveIndicator: {
    flexDirection: "row",
    alignItems: "center",
  },
  liveDot: {
    width: 8,
    height: 8,
    borderRadius: 4,
    backgroundColor: "#8B5CF6",
    marginRight: 6,
  },
  liveText: {
    fontSize: 12,
    fontWeight: "bold",
    color: "#8B5CF6",
  },
  mapContainer: {
    borderRadius: 12,
    overflow: "hidden",
  },
  map: {
    height: 300,
  },
});
