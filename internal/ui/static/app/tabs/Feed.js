import { TvFeedTile } from "../components/TvFeedTile.js";
import { useFeedStore } from "../stores/feed.js";
import { ENDPOINT, PAGE, theme } from "../utils.js";
import { onMounted, ref, storeToRefs } from "../vue.js";

const Feed = {
    props: {

    },
    components: {
        TvFeedTile
    },
    setup: (props) => {

        onMounted(() => {});

        const store = useFeedStore()
        const { feed, loading } = storeToRefs(store)

        return {
            feed,
            loading
        }
    },
    template: `
    <div class="container">
        <div v-if="loading" class="d-flex justify-content-center align-items-center" style="min-height: 50vh;">
            <div class="spinner-border" role="status">
                <span class="visually-hidden">Loading...</span>
            </div>
        </div>
        <div v-else>
            <div v-for="tv in feed" :key="tv.ID" class="mb-3">
                <TvFeedTile :tv="tv"></TvFeedTile>
            </div>
        </div>
    </div>
    `
}
export { Feed };

