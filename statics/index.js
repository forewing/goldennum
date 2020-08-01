var app = new Vue({
    el: '#app',
    data: {

    },
    methods: {
        updateDashboardRoomId(id) {
            this.$refs.panelDashboard.updateRoomId(id);
        },
    }
})