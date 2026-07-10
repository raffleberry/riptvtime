import { currentPage, isMobile, PAGE, theme, updatePage, routes } from "../utils.js";
import { ref } from "../vue.js";


const Navigation = {
    props: {

    },
    components: {

    },
    setup: (props) => {

        const toggleTheme = () => {
            theme.value = theme.value === "light" ? "dark" : "light";
        };

        return {
            c: currentPage,
            theme,
            toggleTheme,
            PAGE,
            routes,
            isMobile
        }
    },
    template: `
    <div class="nav nav-tabs d-flex justify-content-between bg-body" >
      <div class="d-flex flex-row flex-wrap">
        <li v-for="route in routes" :key="route.path" class="nav-item">
            <router-link :class="{'nav-link': true, active: route.path === c.path }" :to="route.path" >{{ route.name }}</router-link>
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

