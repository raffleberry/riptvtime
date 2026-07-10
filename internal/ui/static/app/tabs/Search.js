import { TvSeriesTile } from "../components/TvSeriesTile.js";
import { useSeriesOpts } from "../stores/overlays.js";
import { useSearchStore } from "../stores/search.js";
import { currentPage, ENDPOINT, PAGE, theme, updatePage } from "../utils.js";
import { onMounted, ref, storeToRefs } from "../vue.js";

const Search = {
    props: {
    },
    components: {
        TvSeriesTile
    },
    setup: (props) => {
        
        const store = useSearchStore()

        const { loading, searchTerm, pageCur, results, resultsCnt } = storeToRefs(store)

        const { selected } = storeToRefs(useSeriesOpts())

        onMounted(() => {
            updatePage(PAGE.SEARCH);

            (new bootstrap.Offcanvas(document.querySelector("#seriesOptions"))).show()
            
            selected.value = {
                "Id": 125988,
                "Name": "Silo",
                "Overview": "In a ruined and toxic future, thousands live in a giant silo deep underground. After its sheriff breaks a cardinal rule and residents die mysteriously, engineer Juliette starts to uncover shocking secrets and the truth about the silo.",
                "Status": 1,
                "Year": 2023
            }

        });

      
        return {
            loading,
            searchTerm,
            results,
            resultsCnt,
            pageCur,
        }
    },
    template: `
    <div class="d-flex flex-grow-1 flex-column">
        <div v-if="loading" class="d-flex justify-content-center align-items-center"
            style="min-height: 50vh;">
            <div class="spinner-border" role="status">
                <span class="visually-hidden">Loading...</span>
            </div>
        </div>
        <div v-else-if="resultsCnt === 0" class="d-flex justify-content-center align-items-center" style="min-height: 50vh;">
            <h2>Nothing</h2>
        </div>
        <div v-else class="col">
            <TvSeriesTile class="mb-3" v-for="tv in results[pageCur]" :key="tv.Id" :tv="tv"></TvSeriesTile>
        </div>
    </div>
    `
}
export { Search };

