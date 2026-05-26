(function () {
  const key = "kittycat-theme";
  const root = document.documentElement;

  function apply(theme) {
    root.dataset.theme = theme;
    localStorage.setItem(key, theme);
  }

  const saved = localStorage.getItem(key);
  apply(saved || (matchMedia("(prefers-color-scheme: dark)").matches ? "dark" : "light"));

  document.addEventListener("DOMContentLoaded", () => {
    const btn = document.getElementById("theme-toggle");
    if (!btn) return;
    btn.addEventListener("click", () => {
      apply(root.dataset.theme === "dark" ? "light" : "dark");
    });
  });
})();
