const TvSeriesTile = {
    components: {},
    props: {
        tv: Object
    },
    setup(props) {
        let tv = props.tv
        return {
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
                <h5 class="card-title">{{ tv.Name }} <span class="text-muted">({{ tv.Year }})</span> <span class="badge bg-secondary"</h5>
                <p class="card-text">{{ tv.Overview }}</p>
                </div>
            </div>
        </div>
    </div>
`
}

export { TvSeriesTile };
