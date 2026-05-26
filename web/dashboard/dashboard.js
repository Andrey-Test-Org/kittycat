import { renderChart } from "./components/chart.js";

const STATS_URL = "/stats";

async function loadStats() {
  const res = await fetch(STATS_URL);
  if (!res.ok) throw new Error(`stats ${res.status}`);
  return res.json();
}

function renderCards(stats) {
  const root = document.getElementById("cards");
  root.innerHTML = "";
  for (const [catId, s] of Object.entries(stats)) {
    const card = document.createElement("article");
    card.innerHTML = `
      <h2>${catId}</h2>
      <p>Events: ${s.total_events}</p>
      <p>Avg weight: ${Math.round(s.avg_weight_g)} g</p>
      <p>Kinds: ${s.kinds.join(", ")}</p>
    `;
    root.appendChild(card);
  }
}

async function refresh() {
  try {
    const stats = await loadStats();
    renderCards(stats);
    renderChart(document.getElementById("weight-chart"), stats);
  } catch (err) {
    console.error("refresh failed", err);
  }
}

document.getElementById("refresh").addEventListener("click", refresh);
refresh();
