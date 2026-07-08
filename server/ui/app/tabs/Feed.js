import { TvTile } from "../components/TvTile.js";
import { currentPage, PAGE, theme, updatePage } from "../utils.js";
import { onMounted, ref } from "../vue.js";

let calledOnce = false;
const loading = ref(true)
const feedData = ref([])

const fetchFeed = async () => {
    let url = `/api/series/feed`;
    loading.value = true
    try {
        const response = await fetch(url);
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
        TvTile
    },
    setup: (props) => {

        onMounted(() => {
            updatePage(PAGE.FEED);
            if (!calledOnce) {
                calledOnce = true;
                fetchFeed();
            }
        });

        let testData = {
            "CreatedAt": "2026-07-07T16:31:23.416745591+05:30",
            "DeletedAt": null,
            "EpisodesAired": 144,
            "EpisodesTotal": 144,
            "EpisodesWatched": 3,
            "FirstAirDate": "2018-10-16T00:00:00Z",
            "Genres": "Crime,Drama,Comedy",
            "ID": 2,
            "Name": "The Rookie",
            "Overview": "Starting over isn't easy, especially for small-town guy John Nolan who, after a life-altering incident, is pursuing his dream of being an LAPD officer. As the force's oldest rookie, he's met with skepticism from some higher-ups who see him as just a walking midlife crisis.",
            "RuntimeApprox": 0,
            "TmdbId": 79744,
            "TrackingStatus": 0,
            "UpNextE": 18,
            "UpNextS": 8,
            "UpToDate": false,
            "UpdatedAt": "2026-07-07T16:31:23.416745591+05:30",
            "Year": 2018
        };


        return {
            testData,
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
                <TvTile :tv="tv"></TvTile>
            </div>
        </div>
    </div>
    `
}
export { Feed };

