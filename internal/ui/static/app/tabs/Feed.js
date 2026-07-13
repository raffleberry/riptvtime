import { TvFeedTile } from "../components/TvFeedTile.js";
import { ENDPOINT, PAGE, theme } from "../utils.js";
import { onMounted, ref } from "../vue.js";

let calledOnce = false;
const loading = ref(true)
const feedData = ref([])
// TODO: Refactor this logic into a store 
const fetchFeed = async () => {
    loading.value = true
    try {
        const response = await fetch(ENDPOINT.FEED());
        const result = await response.json();
        if (result) {
            console.log(result)
            feedData.value = result
        } else {
            console.error('Error server sent bad result:', result);
        }
    } catch (error) {
        console.error('Error fetching feed data:', error);
    } finally {
        loading.value = false
    }
}

const Feed = {
    props: {

    },
    components: {
        TvFeedTile
    },
    setup: (props) => {

        onMounted(() => {
            if (!calledOnce) {
                calledOnce = true;
                fetchFeed();
            }
        });

        return {
            feedData,
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
            <div v-for="tv in feedData" :key="tv.ID" class="mb-3">
                <TvFeedTile :tv="tv"></TvFeedTile>
            </div>
        </div>
    </div>
    `
}
export { Feed };

