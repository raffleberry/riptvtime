import { PageRouter } from "./components/PageRouter.js";
import { Confirm } from "./overlays/Confirm.js";
import { Feed } from "./tabs/Feed.js";
import { Upcoming } from "./tabs/Upcoming.js";
import { PAGE, routes } from "./utils.js";
import { computed, createApp, createPinia, createRouter, createWebHistory, onMounted, ref, watch } from "./vue.js";


const router = createRouter({
    history: createWebHistory(),
    routes
});

router.beforeEach((to, from, next) => {
    document.title = to.name || "Tv"
    next()
});

const app = createApp({
    components: {
        PageRouter,
        Confirm,
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
