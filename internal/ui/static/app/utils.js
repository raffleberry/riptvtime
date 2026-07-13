import { Discover } from "./tabs/Discover.js";
import { Feed } from "./tabs/Feed.js";
import { My } from "./tabs/My.js";
import { Search } from "./tabs/Search.js";
import { Series } from "./tabs/Series.js";
import { Stats } from "./tabs/Stats.js";
import { Upcoming } from "./tabs/Upcoming.js";
import { ref, watch } from "./vue.js";

export const ENDPOINT = Object.freeze({
    FEED: ()=>{return '/api/series/feed'},
    SEARCH_SERIES: ()=>{return '/api/series/search'},
    SERIES_STATUS: (Id)=>{return `/api/series/${Id}/status`},
    SERIES_ADD: ()=>{return `/api/series`},
    SERIES_REM: (Id)=>{return `/api/series/${Id}`},
    SERIES_GET: (Id)=>{return `/api/series/${Id}`},
})

export const PAGE = Object.freeze({
    FEED: { name: 'Feed', path: '/' },
    UPCOMING: { name: 'Upcoming', path: '/upcoming' },
    DISCOVER: { name: 'Discover', path: '/discover' },
    SERIES: { name: 'Series', path: '/series/:id' },
    SEARCH: { name: 'Search', path: '/search' },
    STATS: { name: 'Stats', path: '/stats' },
})

export const TvStatus = Object.freeze({
    NotWatching: 0,
	Watching: 1,
	Stopped: 2,
	Completed: 3,
});

export const routes = Object.freeze([
    { path: PAGE.FEED.path, name: PAGE.FEED.name , component: Feed },
    { path: PAGE.UPCOMING.path, name: PAGE.UPCOMING.name , component: Upcoming },
    { path: PAGE.DISCOVER.path, name: PAGE.DISCOVER.name , component: Discover },
    { path: PAGE.SERIES.path, name: PAGE.SERIES.name , component: Series },
    { path: PAGE.SEARCH.path, name: PAGE.SEARCH.name , component: Search },
    { path: PAGE.STATS.path, name: PAGE.STATS.name , component: Stats },
]);

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