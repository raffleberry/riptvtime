import { ENDPOINT } from "../utils.js";
import { defineStore, ref } from "../vue.js";

export const useFeedStore = defineStore('feed', () => {
    const loading = ref(false)
    const feed = ref([])

    const fetchFeed = async () => {
        loading.value = true
        try {
            const response = await fetch(ENDPOINT.FEED());
            if (response.status === 200) {
                const result = await response.json();
                feed.value = result
            } else {
                const msg =`${response.status} - ${await response.text()}`
                throw new Error(msg)
            }
        } catch (error) {
            console.error('Error fetching feed data:', error);
        } finally {
            loading.value = false
        }
    }

    fetchFeed()

    return {
        // data
        loading,
        feed,

        // actions

    }

})