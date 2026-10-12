import { isMobile, PAGE, routes, theme } from "../utils.js"
import { computed, useRoute } from "../vue.js"

const NAV_ICONS = Object.freeze({
  Feed: "bi-house-fill",
  Upcoming: "bi-calendar-event",
  Discover: "bi-compass",
  Search: "bi-search",
  Stats: "bi-bar-chart",
  Import: "bi-upload",
})

const Navigation = {
  props: {},
  components: {},
  setup: (props) => {
    const r = useRoute()

    const curPath = computed(() => r.path)

    const toggleTheme = () => {
      theme.value = theme.value === "light" ? "dark" : "light"
    }

    const tabs = routes.filter((r) => ![PAGE.SERIES.path, PAGE.UNRESOLVED.path].includes(r.path))

    const isActive = (tab) => {
      if (tab.path === curPath.value) return true
      // Discover view-all (?view=...) still belongs to Discover.
      if (tab.path === PAGE.DISCOVER.path && curPath.value === PAGE.DISCOVER.path) return true
      return false
    }

    return {
      curPath,
      theme,
      toggleTheme,
      PAGE,
      tabs,
      isMobile,
      isActive,
      NAV_ICONS,
    }
  },
  template: /* HTML */ `
    <div>
      <!-- Desktop / tablet top nav -->
      <div
        class="mx-2 nav nav-pills py-2 nav-fill d-none d-md-flex justify-content-between bg-body"
      >
        <div class="d-flex flex-row flex-wrap align-items-center">
          <li><h2 class="mb-0 me-2">🪦riptvtime🥀🕊️</h2></li>
          <li v-for="tab in tabs" :key="tab.path" class="nav-item">
            <router-link
              class="nav-link"
              :class="{ active: isActive(tab) }"
              :aria-current="isActive(tab) ? 'page' : null"
              :to="tab.path"
              >{{ tab.name }}</router-link
            >
          </li>
        </div>
        <div class="d-flex flex-row-reverse align-items-center">
          <button class="btn btn-link" @click="toggleTheme" aria-label="Toggle theme">
            <i
              :class="{'bi bi-brightness-high-fill': theme !== 'light', 'bi bi-moon-fill': theme !== 'dark'}"
              style="font-size: 1.25rem;"
            ></i>
          </button>
        </div>
      </div>
      <!-- Mobile compact top bar -->
      <div class="d-flex d-md-none align-items-center justify-content-between px-3 py-2 bg-body">
        <span class="fs-5 fw-bold text-truncate">🪦riptvtime</span>
        <button class="btn btn-link p-2" @click="toggleTheme" aria-label="Toggle theme">
          <i
            :class="{'bi bi-brightness-high-fill': theme !== 'light', 'bi bi-moon-fill': theme !== 'dark'}"
            style="font-size: 1.25rem;"
          ></i>
        </button>
      </div>
      <!-- Mobile bottom tab bar -->
      <nav class="rtt-bottom-nav d-md-none fixed-bottom bg-body border-top" aria-label="Primary">
        <div class="d-flex justify-content-around align-items-stretch h-100 px-1">
          <router-link
            v-for="tab in tabs"
            :key="tab.path"
            class="nav-link flex-fill text-center"
            :class="{ active: isActive(tab) }"
            :aria-current="isActive(tab) ? 'page' : null"
            :to="tab.path"
            :aria-label="tab.name"
          >
            <i :class="'bi ' + (NAV_ICONS[tab.name] || 'bi-circle')"></i>
            <span>{{ tab.name }}</span>
          </router-link>
        </div>
      </nav>
    </div>
  `,
}
export { Navigation }
