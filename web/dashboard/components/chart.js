export function renderChart(canvas, stats) {
  const ctx = canvas.getContext("2d");
  const entries = Object.entries(stats);
  ctx.clearRect(0, 0, canvas.width, canvas.height);

  if (entries.length === 0) {
    ctx.fillStyle = "#888";
    ctx.fillText("no data", 20, 20);
    return;
  }

  const max = Math.max(...entries.map(([, s]) => s.avg_weight_g));
  const barW = canvas.width / entries.length;

  entries.forEach(([catId, s], i) => {
    const h = (s.avg_weight_g / max) * (canvas.height - 40);
    ctx.fillStyle = "#d96c4a";
    ctx.fillRect(i * barW + 8, canvas.height - h - 20, barW - 16, h);
    ctx.fillStyle = "#1f1f1f";
    ctx.fillText(catId, i * barW + 12, canvas.height - 4);
  });
}
