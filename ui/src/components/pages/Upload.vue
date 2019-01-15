<template>
    <div class="columns" style="margin-top:100px">
        <div class="column is-8 is-offset-2">
            <tabs type="is-boxed">
                <tab name="File" icon="ion-android-folder-open" :selected="true">
                    <div class="tile is-child box">
                        <transition name="slide" mode="out-in">
                            <notification type="is-danger" 
                                    @closeNotif="close()" 
                                    v-if="notifActive">
                                {{notificationError}}
                            </notification>
                        </transition>
                        <div class="is-centered">
                            <div class="file-container">
                                <div @click="openFile()">
                                    <svg version="1.1" xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink" x="0px" y="0px"
                                        width="70px" height="70px" viewBox="0 0 516.375 516.375" style="enable-background:new 0 0 70 70;"
                                        xml:space="preserve">
                                        <g fill="#00d1b2">
                                            <path fill="#00d1b2" d="M95.625,38.25c0-9.562,7.65-19.125,19.125-19.125H306v95.625c0,21.038,17.213,38.25,38.25,38.25h76.5v133.875h19.125v-153
                                                L325.125,0H114.75C93.712,0,76.5,17.212,76.5,38.25v248.625h19.125V38.25z M325.125,28.688l89.888,105.188H344.25
                                                c-9.562,0-19.125-9.562-19.125-19.125V28.688z"/>
                                                
                                            <polygon fill="#ccc" points="459,267.75 459,306 57.375,306 57.375,267.75 38.25,267.75 38.25,363.375 57.375,363.375 57.375,325.125 
                                                459,325.125 459,363.375 478.125,363.375 478.125,267.75 		"/>
                                            <g fill="#eee">
                                                <path d="M95.625,478.125V459H76.5v19.125c0,21.037,17.212,38.25,38.25,38.25h19.125V497.25H114.75 C105.188,497.25,95.625,489.6,95.625,478.125z"/>
                                                <rect x="76.5" y="344.25" width="19.125" height="38.25"/>
                                                <rect x="76.5" y="401.625" width="19.125" height="38.25"/>
                                                <rect x="153" y="497.25" width="38.25" height="19.125"/>
                                                <rect x="210.375" y="497.25" width="38.25" height="19.125"/>
                                                <rect x="420.75" y="401.625" width="19.125" height="38.25"/>
                                                <rect x="420.75" y="344.25" width="19.125" height="38.25"/>
                                                <rect x="267.75" y="497.25" width="38.25" height="19.125"/>
                                                <rect x="325.125" y="497.25" width="38.25" height="19.125"/>
                                                <path d="M420.75,478.125c0,9.562-7.65,19.125-19.125,19.125H382.5v19.125h19.125c21.037,0,38.25-17.213,38.25-38.25V459H420.75 V478.125z"/>
                                            </g>
                                        </g>
                                    </svg>
                                </div>
                                <div class="field" style="margin-top:20px">
                                    <div class="file is-primary" :class="{'has-name': filename.length > 1}">
                                        <label class="file-label">
                                            <span class="file-cta">
                                                <span class="file-icon">
                                                    <i class="icon ion-upload"></i>
                                                </span>
                                                <span class="file-label">
                                                    Upload a file
                                                </span>
                                            </span>
                                            <span class="file-name" v-if="filename">
                                                {{filename}}
                                            </span>
                                        </label>
                                    </div>
                                </div>
                                <input class="file-input" type="file" @change="changeHandler($event)">
                            </div>
                        </div>

                        <p class="is-centered" style="margin-top:10px;">
                            <small>
                                By using Saferwall you consent to our <router-link to="/">Terms of Service</router-link> and <router-link to="">Privacy Policy</router-link> and allow us to share your submission with the security community. 
                                <router-link to="">Learn more.</router-link>
                            </small>
                        </p>
                    </div>
                </tab>
                <tab name="Url" icon="ion-ios-world-outline">
                    <div class="tile is-child box">
                        <h1>This part is about our culture</h1>

                        <p class="is-centered" style="margin-top:10px;">
                            <small>
                                By using Saferwall you consent to our <router-link to="/">Terms of Service</router-link> and <router-link to="">Privacy Policy</router-link> and allow us to share your submission with the security community. 
                                <router-link to="">Learn more.</router-link>
                            </small>
                        </p>
                    </div>
                </tab>
            </tabs>
        </div>
    </div>
</template>
<script>
import Tabs from '@/components/elements/Tabs'
import Tab from '@/components/elements/Tab'
import Notification from '@/components/elements/Notification'
import { sha256 } from 'js-sha256';
import axios from 'axios';
import Scanning from '@/components/pages/Scanning';
import Global from '@/global'

export default {
    data(){
        return{
            notificationError: '',
            notifActive: false,
            filename: ''
        }
    },
    components: {
        'tabs': Tabs,
        'tab': Tab,
        'notification': Notification
    },
    methods: {
        changeHandler(e){
            let file = e.target.files[0]
            this.filename = file.name
            if(file.size > 64000000) { 
                this.notificationError = "over sized."
                this.notifActive = true
                return
            }
            
            var reader = new FileReader();
            var hashcode = ''
            reader.onload = (function(theFile) {
                return function(e) {
                    hashcode = sha256(e.target.result)
                    let url = Global.apiUrl + hashcode + '?api-key=' + Global.apiKey
                    axios.get(url, {data: {}}, {headers: {
                            'Content-Type': 'application/json',
                            "Access-Control-Allow-Origin": "*"
                        }})
                        .then((data) => {
                            console.log(data)
                            this.$route.router.go("/scanning")
                        })
                        .catch((err => {
                            console.error(err.response.data)
                        }))
                };
            })(file);
            reader.readAsDataURL(file);
        },
        close(){
            this.notifActive = false
        }
    }
}
</script>
<style lang="scss" scoped>
.slide-enter-active, .slide-leave-active {
  transition: all .3s ease;
}
.slide-enter, .slide-leave-to{
    height: 0;
    padding:0;
    overflow:hidden;
    opacity: 0;
}


.tile{ 
    margin-top: -24px!important; 
    border-radius: 0 0 4px 4px!important; 
    margin-left: 1px!important;
}
.file-container{
    display:inline-block;
    margin: auto;
    position:relative;
    text-align:center;
    cursor:pointer;
    padding: 50px 0;

    input[type="file"]{
        position: absolute;
        top:0;
        left: 0;
        width: 100%;
        height: 100%;
        opacity: 0;
        cursor:pointer;
        z-index: 9;
    }
}
a.button{
    display:inline-block;
    margin-top:10px;
}
</style>

