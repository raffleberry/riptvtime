import { useTracked } from "../../stores/tracked.js";
import { Ky, TvStatus } from "../../utils.js";
import { computed, onMounted, ref, storeToRefs, useRoute, watch } from "../../vue.js";
import { EpisodeOpts } from "./EpisodeOpts.js";
import { SeriesOpts } from "./SeriesOpts.js";
import { useSeriesStore } from "./seriesStore.js";


const Series = {
    props: {
    },
    components: {
        SeriesOpts,
        EpisodeOpts,
    },
    setup: (props) => {
        onMounted(() => {
        });

        const r = useRoute()

        const seriesStore = useSeriesStore()
        const { loading, seriesDetails, SnWatchedEps, epWatchCnt } = storeToRefs(seriesStore)
        const Id = computed(() => r.params.id)

        const { fetchSeries } = seriesStore

        const { series } = storeToRefs(useTracked())

        const { epMarkWatched } = useSeriesStore()

        const seriesId = ref(0)

        const status = computed(() => {
            let ob = series.value?.[Id.value]
            if (ob) {
                return ob.TrackingStatus
            }
            return TvStatus.NotWatching
        })

        const selectedEp = ref({})

        onMounted(() => {
            if (r.hash) {
                let sel = document.querySelector(r.hash)
                if (sel) sel.scrollIntoView()
            }

        })

        const getStatusTxt = (s) => {
            switch (s) {
                case TvStatus.Watching:
                    return "Watching"
                case TvStatus.Completed:
                    return "Completed"
                case TvStatus.Stopped:
                    return "Stopped"
                default:
                    return "Add";
            }
        }

        const statusCss = computed(() => {
            let card = 'position-sticky top-0 card'
            let icon = 'bi'
            let btn = ''
            let pgbar = ''
            switch (status.value) {
                case TvStatus.Watching:
                    card += ' border border-warning'
                    icon += ' bi-bookmark-check'
                    btn += ' btn-outline-warning'
                    pgbar += 'bg-warning'
                    break;
                case TvStatus.NotWatching:
                    // cb += ' border border-primary'
                    icon += ' bi-bookmark-plus'
                    btn += ' btn-outline-primary'
                    pgbar += 'bg-primary'
                    break;
                case TvStatus.Stopped:
                    card += ' border border-danger'
                    icon += ' bi-bookmark-x'
                    btn += ' btn-outline-danger'
                    pgbar += 'bg-danger'
                    break;
                case TvStatus.Completed:
                    card += ' border border-success'
                    icon += ' bi-bookmark-check-fill'
                    btn += ' btn-outline-success'
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

        const progress = computed(() => {
            if (seriesDetails.value?.NumberOfEpisodes) {
                const deno = seriesDetails.value.NumberOfEpisodes
                if (deno === 0) {
                    return 0
                }
                const num = Object.keys(epWatchCnt.value).length
                return (num / deno) * 100
            }

            return 0
        })

        watch(Id, (id) => {
            if (id) {
                fetchSeries(id)
            }
        }, { immediate: true })
        
        var epOptEl = null

        onMounted(() => {
            const el = document.getElementById('episodeOpts')
            epOptEl = bootstrap.Offcanvas.getOrCreateInstance(el)

            // el.addEventListener('show.bs.offcanvas', () => {
            // })

            // el.addEventListener('hidden.bs.offcanvas', () => {
            // })
        })


        const cnt = (s, e) => {
            return epWatchCnt.value[Ky(s, e)] ?? 0
        }

        const openEpOpts = async (ep) => {
            if (cnt(ep.SeasonNumber, ep.EpisodeNumber) === 0) {
                await epMarkWatched(Number(Id.value), ep.SeasonNumber, ep.EpisodeNumber)
            } else {
                selectedEp.value = ep
                epOptEl.show()
            }
        }

        const isAired = (dStr) => {
            return (new Date()) >= (new Date(dStr))
        }


        return {
            Ky,
            loading,
            sd: seriesDetails,
            SnWatchedEps,
            epWatchCnt,
            statusCss,
            status,
            getStatusTxt,
            r,
            openEpOpts,
            selectedEp,
            cnt,
            progress,
            isAired,
        }
    },
    template: `
    <EpisodeOpts :mid="sd.Id" :ep="selectedEp"></EpisodeOpts>
    <SeriesOpts :mid="sd.Id" :name="sd.Name" :year="sd.Year"></SeriesOpts>
    <div class="container-fluid">
        <div v-if="loading" class="d-flex justify-content-center align-items-center"
            style="min-height: 50vh;">
            <div class="spinner-border" role="status">
                <span class="visually-hidden">Loading...</span>
            </div>
        </div>
        <div v-else>
            <div :class="statusCss.card" style="z-index: 5;">
                <div class="card-body">
                    <div class="d-flex justify-content-between">
                        <h3 class="card-title">{{ sd.Name }} <span class="text-muted">({{ sd.Year }}) </span></h3>
                        <button :class="statusCss.btn"
                            data-bs-toggle="offcanvas" data-bs-target="#seriesOpts"
                            type="button" class="btn btn-sm p-2 d-inline-flex align-items-center justify-content-center">
                                <i class="bi me-2" :class="statusCss.icon"></i>
                                {{ getStatusTxt(status) }}
                        </button>
                    </div>
                    <p class="card-text">{{ sd.Overview }}</p>
                </div>
                <div>
                    <div class="progress" role="progressbar" aria-label="Tv show progress" :aria-valuenow="progress" aria-valuemin="0" aria-valuemax="100">
                        <div :class="statusCss.pgbar" class="progress-bar" :style="{width: progress + '%'}"></div>
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
                    <div
                        :class="{ 'show': '#season' + sn.SeasonNumber === r.hash }"
                        :id="'season' + sn.SeasonNumber" class="accordion-collapse collapse"
                        data-bs-parent="#seasons">
                        <div class="accordion-body">
                            <div v-for="(ep, idx) in sn.Episodes" 
                                :class="idx % 2 === 0 ? 'bg-body' : 'bg-body-secondary'"
                                class="d-flex flex-row justify-content-between align-items-center">
                                <div>{{ ep.SeasonNumber }}x{{ ep.EpisodeNumber }} - {{ ep.Name }}</div>
                                <div>
                                    <span v-if="cnt(ep.SeasonNumber, ep.EpisodeNumber) > 1">{{ cnt(ep.SeasonNumber, ep.EpisodeNumber) }}x</span>
                                    <span v-if="!isAired(ep.AirDate)">{{ new Date(ep.AirDate).toDateString() }}</span>
                                    <button :disabled="!isAired(ep.AirDate)" class="btn" :class="!isAired(ep.AirDate) ? 'border-0' : ''"
                                        @click="() => { openEpOpts(ep) }"
                                    >
                                        <i class="bi" :class="cnt(ep.SeasonNumber, ep.EpisodeNumber) > 0 ? 'bi-check-circle-fill text-success' : 'bi-check-circle'"></i>
                                    </button>
                                </div>
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

