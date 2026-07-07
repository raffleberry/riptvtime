import { currentPage, isMobile, PAGE, theme, updatePage } from "../utils.js";
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

        const isFullscreen = ref(false)

        const toggleFullscreen = async () => {
            if (document.fullscreenElement) {
                await document.exitFullscreen();
            } else {
                await document.documentElement.requestFullscreen();
            }
            isFullscreen.value = document.fullscreenElement !== null
        }

        return {
            c: currentPage,
            theme,
            toggleTheme,
            PAGE,
            toggleFullscreen,
            isFullscreen,
            isMobile
        }
    },
    template: `
    <div class="nav nav-tabs d-flex justify-content-between bg-body" >
      <div class="d-flex flex-row flex-wrap">
        <li class="nav-item"><router-link :class="{'nav-link': true, active: c.name === PAGE.FEED.name }" :to="PAGE.FEED.path" >Feed</router-link></li>
        <li class="nav-item"><router-link :class="{'nav-link': true, active: c.name === PAGE.UPCOMING.name }" :to="PAGE.UPCOMING.path" >Upcoming</router-link></li>
      </div>
      <div class="d-flex flex-grow-1 flex-row-reverse">
        <button v-if="isMobile" class="btn btn-link" @click="toggleFullscreen" aria-label="fullscreenToggle">
            <i :class="[ isFullscreen ? 'bi-fullscreen-exit' : 'bi-arrows-fullscreen' ]" style="font-size: 1.25rem;"></i>
        </button>
        <button class="btn btn-link" @click="toggleTheme" aria-label="darkMode">
          <i :class="{'bi bi-brightness-high-fill': theme !== 'light', 'bi bi-moon-fill': theme !== 'dark'}" style="font-size: 1.25rem;"></i>
        </button>
      </div>

    </div>
    `
}
export { Navigation };

