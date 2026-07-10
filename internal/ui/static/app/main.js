import { PageRouter } from "./components/PageRouter.js";
import { SeriesOptions } from "./overlays/SeriesOptions.js";
import { Feed } from "./tabs/Feed.js";
import { Upcoming } from "./tabs/Upcoming.js";
import { PAGE, currentPage, routes } from "./utils.js";
import { computed, createApp, createPinia, createRouter, createWebHistory, onMounted, ref, watch } from "./vue.js";


const router = createRouter({
    history: createWebHistory(),
    routes
});

const app = createApp({
    components: {
        PageRouter,
        SeriesOptions,
    },

    setup() {

        onMounted(() => {
        });

        return {
            
        };
    }
})

app.use(router)
app.use(createPinia())
app.mount('#app')
