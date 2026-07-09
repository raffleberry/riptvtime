import { TvSeriesTile } from "../components/TvSeriesTile.js";
import { currentPage, ENDPOINT, PAGE, theme, updatePage } from "../utils.js";
import { onMounted, ref } from "../vue.js";

export const searchLoading = ref(false)

let searchTerm = ''

export const searchResults = ref({
    Page: 1,
    Results: {
        1: [],
    },
    TotalPages: 0,
    TotalResults: 0,
})

export const handleSearch = async (searchText) => {
    
    if (searchLoading.value || searchTerm === searchText) return

    try {
        searchLoading.value = true
        const { data, err } = await searchTv(searchText, 1)
        if (err) {
            throw err
        }

        searchTerm = searchText
        searchResults.value = {
            Page: 1,
            Results: {
                1: data.Results,
            },
            TotalPages: data.TotalPages,
            TotalResults: data.TotalResults,
        }
        console.log(searchResults.value, data)
    } catch (error) {
        console.error(error)
    } finally {
        searchLoading.value = false
    }

}

export const handleSearchNxtBtn = async () => {

    if (searchResults.value.Results[searchResults.value.Page + 1]) {
        searchResults.value.Page += 1
        return
    }
    
    
    try {
        let curPg = searchResults.value.Page
        searchLoading.value = true
        const { data, err } = await searchTv(searchTerm, curPg + 1)
        if (err) {
            throw err
        }


        searchResults.value = {
            Page: curPg + 1,
            Results: {
                ...searchResults.value.Results,
                [curPg + 1]: data.Results,
            },
            TotalPages: data.TotalPages,
            TotalResults: data.TotalResults,
        }

    } catch (error) {
        console.error(error)
    } finally {
        searchLoading.value = false
    }

}

export const handleSearchPrvBtn = async () => {

    if (searchResults.value.Results[searchResults.value.Page - 1]) {
        searchResults.value.Page -= 1
        return
    }
    
    
    try {
        let curPg = searchResults.value.Page
        searchLoading.value = true
        const { data, err } = await searchTv(searchTerm, curPg - 1)
        if (err) {
            throw err
        }


        searchResults.value = {
            Page: curPg - 1,
            Results: {
                ...searchResults.value.Results,
                [curPg - 1]: data.Results,
            },
            TotalPages: data.TotalPages,
            TotalResults: data.TotalResults,
        }

    } catch (error) {
        console.error(error)
    } finally {
        searchLoading.value = false
    }

}


const searchTv = async (search, page) => {
    let url = `${ENDPOINT.SEARCH_SERIES}?q=${search}&p=${page}`
    try {
        const response = await fetch(url);
        const result = await response.json();
        return {
            data: result,
            err: null
        }
    } catch (error) {
        console.error('Error fetching music data:', error);
        return {
            data: result,
            err: error
        }
    }
}

const Search = {
    props: {
    },
    components: {
        TvSeriesTile
    },
    setup: (props) => {

        onMounted(() => {
            updatePage(PAGE.SEARCH);
        });

        const onSearch = async () => {
            searchLoading.value = true
            let result = await searchTv(searchTerm, 1)
            if (result.data) {
                searchResults.value = result.data
            }
            searchLoading.value = false
        }


        return {
            loading: searchLoading,
            searchTerm,
            searchResults,
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
        <div v-else-if="searchResults.Results.TotalResults === 0" class="d-flex justify-content-center align-items-center" style="min-height: 50vh;">
            No results
        </div>
        <div v-else class="col">
            <TvSeriesTile class="mb-3" v-for="tv in searchResults.Results[searchResults.Page]" :key="tv.Id" :tv="tv"></TvSeriesTile>
        </div>
    </div>
    `
}
export { Search };

