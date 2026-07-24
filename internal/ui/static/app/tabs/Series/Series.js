import { useConfirm } from "../../stores/confirm.js";
import { useTracked } from "../../stores/tracked.js";
import { Ky, PAGE, theme, TvStatus } from "../../utils.js";
import { computed, onMounted, ref, storeToRefs, useRoute, watch } from "../../vue.js";
import { useSeriesStore } from "./seriesStore.js";


const Series = {
    props: {

    },
    components: {
    },
    setup: (props) => {
        onMounted(() => {
        });

        const r = useRoute()
        const seriesStore = useSeriesStore()
        const { loading, seriesDetails, SnWatchedEps, EpWatchCnt } = storeToRefs(seriesStore)
        const Id = computed(() => r.params.id)

        const confirmStore = useConfirm()
        const { openDialog } = confirmStore

        const { fetchSeries } = seriesStore

        const { series } = storeToRefs(useTracked())

        const seriesId = ref(0)

        const status = computed(() => {
            let ob = series.value?.[Id.value]
            if (ob) {
                return ob.TrackingStatus
            }
            return TvStatus.NotWatching
        })

        const statusCss = computed(() => {
            let card = 'card'
            let icon = 'bi'
            let btn = 'btn'
            let pgbar = ''
            switch (status.value) {
                case TvStatus.Watching:
                    card += ' border border-warning'
                    icon += ' bi-bookmark-check'
                    btn += ' btn-warning'
                    pgbar += 'bg-warning'
                    break;
                case TvStatus.NotWatching:
                    // cb += ' border border-primary'
                    icon += ' bi-three-dots-vertical'
                    btn += ' btn-primary'
                    pgbar += 'bg-primary'
                    break;
                case TvStatus.Stopped:
                    card += ' border border-danger'
                    icon += ' bi-bookmark-x'
                    btn += ' btn-danger'
                    pgbar += 'bg-danger'
                    break;
                case TvStatus.Completed:
                    card += ' border border-success'
                    icon += ' bi-bookmark-check-fill'
                    btn += ' btn-success'
                    pgbar += 'bg-success'
                    break;

                default:
                    break;
            }
            return {
                card,
                icon,
                btn,
                pgbar,
            };
        });

        watch(Id, (id) => {
            if (id) {
                fetchSeries(id)
            }
        }, { immediate: true })

        return {
            Ky,
            loading,
            sd: seriesDetails,
            SnWatchedEps,
            EpWatchCnt,
            openDialog,
            statusCss,
        }
    },
    template: `
    <div class="container-fluid">
        <div v-if="loading" class="d-flex justify-content-center align-items-center"
            style="min-height: 50vh;">
            <div class="spinner-border" role="status">
                <span class="visually-hidden">Loading...</span>
            </div>
        </div>
        <div v-else>
            <div :class="statusCss.card">
                <div class="card-body">
                    <div class="d-flex justify-content-between">
                        <h3 class="card-title">{{ sd.Name }} <span class="text-muted">({{ sd.Year }}) </span></h3>
                        <button :class="statusCss.btn"
                            @click="openSeriesOptions" 
                            data-bs-toggle="offcanvas" data-bs-target="#seriesOptions"
                            type="button" class="btn p-2 d-inline-flex align-items-center justify-content-center">
                                <i :class="statusCss.icon">Options</i>
                        </button>
                    </div>
                    <p class="card-text">{{ sd.Overview }}</p>
                </div>
                <div>
                    <div class="progress" role="progressbar" aria-label="Tv show progress" :aria-valuenow="watchProgress" aria-valuemin="0" aria-valuemax="100">
                        <div class="progress-bar" :style="{width: watchProgress + '%'}"></div>
                    </div>
                </div>
            </div>
            <div class="accordion" id="seasons">

                <div v-for="sn in sd.Seasons" class="accordion-item">
                    <h2 class="accordion-header">
                        <button class="accordion-button" type="button" data-bs-toggle="collapse" :data-bs-target="'#season' + sn.SeasonNumber">
                            <div class="d-flex w-100 justify-content-between pe-5"><span>{{ sn.Name }}</span> <span>{{SnWatchedEps[sn.SeasonNumber]?.length}}/{{sn.EpisodeCount}}</span></div>
                        </button>
                    </h2>
                    <div :id="'season' + sn.SeasonNumber" class="accordion-collapse collapse" data-bs-parent="#seasons">
                        <div class="accordion-body">
                            <div class="d-flex flex-row justify-content-between" v-for="ep in sn.Episodes">
                                <div>{{ ep.SeasonNumber }}x{{ ep.EpisodeNumber }} - {{ ep.Name }}</div>
                                <button class="btn"
                                    @click="() => { openDialog(Ky(ep.SeasonNumber, ep.EpisodeNumber)) }"
                                >
                                    <i :class="Ky(ep.SeasonNumber, ep.EpisodeNumber) in EpWatchCnt ? 'bi bi-check-circle-fill text-success' : 'bi bi-check-circle'"></i>
                                </button>
                            </div>
                        </div>
                    </div>
                </div>

            </div>
        </div>
    </div>
    `
}
export { Series };

