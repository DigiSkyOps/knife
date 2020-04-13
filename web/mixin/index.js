import { mapActions } from "vuex"
import { Notification } from "element-ui"
export const index = {
    data(){
        return{

        }
    },
    methods: {
        initListen(obs){
            obs.on("ConnectionOpened",()=>{
                console.log(this)

            })
            obs.on('ConnectionClosed',()=>{

            })
            obs.on('error', err => {
                Notification.error({
                    title: "socket error!",
                    message: err
                })
            })
        },
    },
}
