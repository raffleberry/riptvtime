import { isMobile, PAGE, theme, routes } from "../utils.js";
import { computed, ref, useRoute } from "../vue.js";


const Navigation = {
    props: {

    },
    components: {

    },
    setup: (props) => {

        const r = useRoute()

        const curPath = computed(() => r.path)

        const toggleTheme = () => {
            theme.value = theme.value === "light" ? "dark" : "light";
        };

        const tabs = routes.filter(r => ![PAGE.SERIES.path].includes(r.path));

        return {
            curPath,
            theme,
            toggleTheme,
            PAGE,
            tabs,
            isMobile
        }
    },
    template: `
    <div class="nav nav-tabs d-flex justify-content-between bg-body" >
      <div class="d-flex flex-row flex-wrap">
        <li v-for="tab in tabs" :key="tab.path" class="nav-item">
            <router-link :class="{'nav-link': true, active: tab.path === curPath }" :to="tab.path" >{{ tab.name }}</router-link>
        </li>

      </div>
      <div class="d-flex flex-grow-1 flex-row-reverse">
        <button class="btn btn-link" @click="toggleTheme" aria-label="darkMode">
          <i :class="{'bi bi-brightness-high-fill': theme !== 'light', 'bi bi-moon-fill': theme !== 'dark'}" style="font-size: 1.25rem;"></i>
        </button>
      </div>

    </div>
    `
}
export { Navigation };

