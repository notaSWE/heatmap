async function createHeatmap() {
  const rawData = await d3.csv("data.csv");

  const items = rawData
  .filter(row => parseFloat(row.portweight) >= 0.02) // Filter out tickers with portweight under 2%
  .map(row => {
      return {
          ticker: row.ticker,
          value: parseFloat(row.portweight),
          daychange: parseFloat(row.daychange),
      };
  });

  const colorScale = d3.scaleLinear()
      .domain([-3, 0, 3])
      .range(["rgb(246, 53, 56)", "rgb(65, 69, 84)", "rgb(48, 204, 90)"])
      .clamp(true);

      function renderHeatmap() {
        heatmap.innerHTML = ''; // Clear existing rectangles

        const treemapLayout = d3.treemap()
            .size([window.innerWidth, window.innerHeight])
            .padding(1)
            .round(true);

        const root = d3.hierarchy({ children: items })
            .sum(d => d.value)
            .sort((a, b) => b.value - a.value);

        treemapLayout(root);

        root.leaves().forEach(leaf => {
          const rectangle = document.createElement('div');
          rectangle.className = 'rectangle';
          rectangle.style.left = `${leaf.x0}px`;
          rectangle.style.top = `${leaf.y0}px`;
          rectangle.style.width = `${leaf.x1 - leaf.x0}px`;
          rectangle.style.height = `${leaf.y1 - leaf.y0}px`;
          rectangle.style.backgroundColor = colorScale(leaf.data.daychange);
    
          const label = document.createElement('div');
          label.className = 'label';
    
          const tickerText = document.createElement('span');
          tickerText.className = 'ticker';
          tickerText.textContent = `${leaf.data.ticker}`;
    
          const daychangeText = document.createElement('span');
          daychangeText.className = 'daychange';
          daychangeText.textContent = leaf.data.daychange > 0 ? `+${leaf.data.daychange.toFixed(2) + "%"}` : `${leaf.data.daychange.toFixed(2) + "%"}`;
    
          label.appendChild(tickerText);
          label.appendChild(daychangeText);
          rectangle.appendChild(label); // Add the label to the rectangle
    
          heatmap.appendChild(rectangle);
      });
    }

    window.addEventListener('resize', renderHeatmap); // Add an event listener to call renderHeatmap on window resize

    renderHeatmap(); // Call renderHeatmap initially
}

createHeatmap();
