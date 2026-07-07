import { PageRouter } from "./PageRouter.js";
import { Feed } from "./tabs/Feed.js";
import { Upcoming } from "./tabs/Upcoming.js";
import { PAGE, currentPage } from "./utils.js";
import { computed, createApp, createPinia, createRouter, createWebHistory, onMounted, ref, watch } from "./vue.js";


const routes = [
    { path: PAGE.FEED.path, component: Feed, name: PAGE.FEED.name },
    { path: PAGE.UPCOMING.path, component: Upcoming, name: PAGE.UPCOMING.name },
];

const router = createRouter({
    history: createWebHistory(),
    routes
});

const app = createApp({
    components: {
        PageRouter,
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
