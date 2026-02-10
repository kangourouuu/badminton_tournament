import { ref, onMounted, onUnmounted, nextTick } from "vue";

export function useBracketConnectors(containerRef) {
  const paths = ref([]);

  // Calculate center of an element relative to the container
  const getElementPoint = (id, side = "center") => {
    const el = document.getElementById(id);
    if (!el || !containerRef.value) return null;

    const rect = el.getBoundingClientRect();
    const cRect = containerRef.value.getBoundingClientRect();

    const x =
      side === "right"
        ? rect.right - cRect.left
        : side === "left"
          ? rect.left - cRect.left
          : rect.left - cRect.left + rect.width / 2;

    const y = rect.top - cRect.top + rect.height / 2;

    return { x, y };
  };

  const drawOrthogonal = (p1, p2) => {
    // M x1 y1 L xMid y1 L xMid y2 L x2 y2
    const midX = (p1.x + p2.x) / 2;
    return `M ${p1.x} ${p1.y} L ${midX} ${p1.y} L ${midX} ${p2.y} L ${p2.x} ${p2.y}`;
  };

  const drawOrthogonalDashed = (p1, p2) => {
    // Same as orthogonal but meant for styling externally
    return drawOrthogonal(p1, p2);
  };

  // connections: Array of { startId, endId, type: 'orthogonal' | 'straight' | 'dashed' }
  const updateConnectors = async (connections) => {
    await nextTick();
    if (!containerRef.value) return;

    const newPaths = [];

    connections.forEach((conn) => {
      const p1 = getElementPoint(conn.startId, "right"); // Start from right of source
      const p2 = getElementPoint(conn.endId, "left"); // End at left of target

      // For Bronze match, we might want to start from right similarly

      if (p1 && p2) {
        let d = "";
        if (conn.type === "orthogonal" || conn.type === "dashed") {
          d = drawOrthogonal(p1, p2);
        } else {
          // Default straight or whatever
          d = `M ${p1.x} ${p1.y} L ${p2.x} ${p2.y}`;
        }

        newPaths.push({
          d,
          type: conn.type,
          key: `${conn.startId}-${conn.endId}`,
        });
      }
    });

    paths.value = newPaths;
  };

  return {
    paths,
    updateConnectors,
  };
}
