import { Navigation } from "./Navigation.js";
import { onMounted, ref, RouterView } from "../vue.js";




const TvFeedTile = {
    components: {},
    props: {
        tv: Object
    },
    setup(props) {
        console.log(props)
        let tv = props.tv
        let upNext = "S" + String(tv.UpNextS).padStart(2, "0") + "E" + String(tv.UpNextE).padStart(2, "0")
        let toWatchCnt = tv.EpisodesAired - tv.EpisodesWatched - 1
        let watchProgress = (tv.EpisodesWatched / tv.EpisodesAired) * 100
        console.log(tv)
        return {
            upNext,
            toWatchCnt,
            watchProgress,
            tv
        }
    },

    template: `
    <div class="card">
    <div class="row">
        <!--
        <div class="col-md-4">
        <img src="..." class="img-fluid rounded-start" alt="..."> 
        </div>
        <div class="col-md-8">
        </div>
        -->
        <div class="col">
            <div class="card-body">
            <h5 class="card-title">{{ tv.Name }} <span class="text-muted">({{ tv.Year }})</span> <span class="badge bg-secondary" v-if="tv.RecentlyAired">New</span> </h5>
            <p class="card-text">{{ tv.Overview }}</p>
            <p class="card-text">
                <span> Up Next </span>
                <button type="button" class="btn btn-primary position-relative">
                {{ upNext }}
                <span v-if="toWatchCnt > 0" class="position-absolute top-0 start-100 translate-middle badge rounded-pill bg-success">
                    +{{ toWatchCnt }}
                    <span class="visually-hidden">unwatched episodes</span>
                </span>
                </button>
            </p>
            </div>
        </div>
        <div>
            <div class="progress" role="progressbar" aria-label="Tv show progress" :aria-valuenow="watchProgress" aria-valuemin="0" aria-valuemax="100">
                <div class="progress-bar" :style="{width: watchProgress + '%'}"></div>
            </div>
        </div>
    </div>
    </div>
`
}

export { TvFeedTile };
