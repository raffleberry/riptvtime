import { SearchButtons } from "../../components/SearchButtons.js";
import { ENDPOINT, PAGE, theme } from "../../utils.js";
import { computed, onMounted, ref, storeToRefs, useRoute } from "../../vue.js";
import { SearchBox } from "./SearchBox.js";
import { SearchTile } from "./SearchTile.js";
import { SearchTileOpts } from "./SearchTileOpts.js";
import { useSeriesOpts } from "./optsStore.js";
import { useSearchStore } from "./searchStore.js";

const Search = {
    props: {
        
    },
    components: {
        SearchTileOpts,
        SearchTile,
        SearchBox,
        SearchButtons,
    },
    setup: (props) => {
        
        const store = useSearchStore()

        const { loading, searchTerm, pageCur, results, resultsCnt } = storeToRefs(store)

        const { selected } = storeToRefs(useSeriesOpts())


        const r = useRoute()
        const curPath = computed(() => r.path)


        onMounted(() => {

        });

      
        return {
            loading,
            searchTerm,
            results,
            resultsCnt,
            pageCur,
            curPath,
            PAGE,
        }
    },
    template: `
    <SearchTileOpts></SearchTileOpts>
    <SearchBox class="my-2" v-if="curPath === PAGE.SEARCH.path">
    </SearchBox>
    <div class="flex-grow-1 d-flex flex-column overflow-auto">
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
                <SearchTile class="mb-3" v-for="tv in results[pageCur]" :key="tv.Id" :tv="tv"></SearchTile>
            </div>
        </div>
    </div>
    <SearchButtons class="my-2"
        v-if="curPath === PAGE.SEARCH.path && resultsCnt > 0">
    </SearchButtons>
    `
}
export { Search };

