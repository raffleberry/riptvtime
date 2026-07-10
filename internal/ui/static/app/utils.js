import { Feed } from "./tabs/Feed.js";
import { Search } from "./tabs/Search.js";
import { Upcoming } from "./tabs/Upcoming.js";
import { ref, watch } from "./vue.js";

export const ENDPOINT = Object.freeze({
    FEED: ()=>{return '/api/series/feed'},
    SEARCH_SERIES: ()=>{return '/api/series/search'},
    SERIES_STATUS: (mId)=>{return `/api/series/${mId}/status`},
})

export const PAGE = Object.freeze({
    FEED: { name: 'Feed', path: '/' },
    UPCOMING: { name: 'Upcoming', path: '/upcoming' },
    STATS: { name: 'Stats', path: '/stats' },
    MY_SHOWS: { name: 'My Shows', path: '/my' },
    DISCOVER: { name: 'Discover', path: '/discover' },
    SEARCH: { name: 'Search', path: '/search' },
})

export const TvStatus = Object.freeze({
    NotWatching: 0,
	Watching: 1,
	Stopped: 2,
	Completed: 3,
});

// use updatePageTitle to update the page title
export const currentPage = ref(PAGE.FEED)

export const routes = Object.freeze([
    { path: PAGE.FEED.path, component: Feed, name: PAGE.FEED.name },
    { path: PAGE.UPCOMING.path, component: Upcoming, name: PAGE.UPCOMING.name },
    { path: PAGE.SEARCH.path, component: Search, name: PAGE.SEARCH.name },
]);

export const updatePage = (page) => {
    currentPage.value = page;
    document.title = page.name;
}

export const highlight = (text, highlight) => {
    if (!highlight) return text;
    const regex = new RegExp(`(${highlight})`, 'gi');
    return text.replace(regex, `<strong>$1</strong>`);
}

export const formatDuration = (seconds) => {
    const minutes = Math.floor(seconds / 60);
    const secs = Math.floor(seconds % 60);
    return `${minutes}:${secs < 10 ? '0' : ''}${secs}`;
};

const mediaQuery = window.matchMedia('(max-width: 1000px)')
export const isMobile = ref(mediaQuery.matches)
const updateIsMobile = (event) => { isMobile.value = event.matches }
mediaQuery.addEventListener('change', updateIsMobile)

export const theme = ref(localStorage.getItem("theme") || "light");
watch(theme, () => {
    document.documentElement.setAttribute("data-bs-theme", theme.value);
    localStorage.setItem("theme", theme.value);
}, { immediate: true });

export function generateRandomString(length) {
    let result = '';
    const characters = 'ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789';
    const charactersLength = characters.length;
    for (let i = 0; i < length; i++) {
        result += characters.charAt(Math.floor(Math.random() * charactersLength));
    }
    return result;
}