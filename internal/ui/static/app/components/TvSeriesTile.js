import { useSeriesOpts } from "../stores/overlays.js";
import { TvStatus } from "../utils.js";
import { storeToRefs } from "../vue.js";

const TvSeriesTile = {
    components: {},
    props: {
        tv: Object
    },
    setup(props) {
        let tv = props.tv
        const store = useSeriesOpts()
        const { selected } = storeToRefs(store)
        console.log(tv)

        const openSeriesOptions = () => {
            selected.value = {
                MId: tv.Id,
                Name: tv.Name,
                Year: tv.Year,
                Status: tv.Status
            }
        }

        return {
            tv,
            TvStatus,
            openSeriesOptions,
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
                <div class="d-flex justify-content-between">
                    <h5 class="card-title">{{ tv.Name }} <span class="text-muted">({{ tv.Year }})</span> <span class="badge bg-secondary" </span></h5>
                    <button @click="openSeriesOptions" 
                        data-bs-toggle="offcanvas" data-bs-target="#seriesOptions"
                        type="button" class="btn p-2 d-inline-flex align-items-center justify-content-center">
                            <i v-if="tv.Status === TvStatus.Watching" class="bi bi-bookmark-check text-primary"></i>
                            <i v-else-if="tv.Status === TvStatus.NotWatching" class="bi bi-three-dots-vertical"></i>
                            <i v-else-if="tv.Status === TvStatus.Stopped" class="bi bi-bookmark-x text-danger"></i>
                            <i v-else-if="tv.Status === TvStatus.Completed" class="bi bi-bookmark-check-fill text-success"></i>
                        </button>
                </div>
                <p class="card-text">{{ tv.Overview }}</p>
                </div>
            </div>
        </div>
    </div>
`
}

export { TvSeriesTile };
