<template>
  <div class="dashboard">
    <h1>Sin miedo al exito</h1>
    <div v-if="!loading">
    <v-container class="grey lighten-5">
      <v-dialog
      v-model="dialogerror"
      persistent
      max-width="290"
      >
        
        <v-card>
          <v-card-title class="text-h5">
            Date Error
          </v-card-title>
          <v-card-text >error: {{errordetail}}</v-card-text>
          <v-card-actions>
            <v-spacer></v-spacer>
            <v-btn
              color="green darken-1"
              text
              @click="dialogerror = false"
            >
              Close
            </v-btn>
          </v-card-actions>
        </v-card>
      </v-dialog>     
      <v-row >
        <v-col cols ="6" >
          <v-row>
          <v-col cols=3>
            <Date v-on:date="getDate($event)" />
            
          </v-col>
          <v-col>
            <v-btn
              color="primary"
              elevation="2"
              @click="fetchClientsbyDate(date);"
            >Search
            <v-icon right>search</v-icon>
            </v-btn>
          </v-col>
          </v-row>
        
        <v-toolbar
            class="mb-2"
            color="indigo darken-5"
            dark
            flat
          >
            <v-toolbar-title>Clients List</v-toolbar-title>
          </v-toolbar>
          <v-card>
            <v-card-title>
              <v-text-field
                v-model="search"
                append-icon="mdi-magnify"
                label="Search"
                single-line
                hide-details
              ></v-text-field>
            </v-card-title>
              <v-data-table
                :search="search"
                :headers="headers"
                :items="datos"
                :items-per-page="5"
                class="elevation-1"
                @click:row="handleclick"            
              >        
            </v-data-table>
          </v-card>

        </v-col>
      </v-row>
      <v-row>
        <div v-if="clientRecords">
        <v-col >
          <Popup  v-bind:Tdata="Tdata" v-bind:similarBuyer="similarBuyer" v-bind:clientName="clientName"
          v-bind:favoriteProducts="favoriteProducts" />
        </v-col>
        </div>
        <div v-else-if="clientHasNoRecords">
        <v-col >
          <h2>Cant find any records for that user</h2>
        </v-col>
        </div>
        <div v-else>          
          <h2> Fetching client data</h2>
        </div>
      </v-row>      
    </v-container>
    </div>
    <div v-else >
      <h2> Loading data....</h2>
      <content-loader
      viewBox="0 0 476 124"
      primaryColor="#f3f3f3"
      secondaryColor="#cccccc"
    >
      <rect x="48" y="8" rx="3" ry="3" width="88" height="6" />
      <rect x="48" y="26" rx="3" ry="3" width="52" height="6" />
      <rect x="0" y="56" rx="3" ry="3" width="410" height="6" />
      <rect x="0" y="72" rx="3" ry="3" width="380" height="6" />
      <rect x="0" y="88" rx="3" ry="3" width="178" height="6" />
      <circle cx="20" cy="20" r="20" />
    </content-loader>
    
     
    </div>
  </div>
</template>

<script>
import axios from 'axios'
import Date from '@/components/Date.vue'
import Popup from '../components/Popup.vue'
import { ContentLoader } from "vue-content-loader";



export default {  
  data(){
    return{
      ok:true,
      errordetail:null,
      clientRecords:false,
      dialogerror:false,
      clientHasNoRecords:false,
      loading:true,
      clientName:'',
      date:null,
      search: '',
      currentClient:{
            Cid: '',
            name: "",
            age: '',
        

      },
      datos:null,
      Tdata:null,
      Theader:[
        { text: 'Transaction No', value: 'Tid' },
        { text: 'Products', value: 'ProductIds' },
        { text: 'TOTAL', value: 'price' },


      ],
      similarBuyer:null,
      favoriteProducts:null,

      headers:[
        {
            text: 'Client ID',
            align: 'start',
            sortable: false,
            value: 'Cid',
          },
          { text: 'Name', value: 'name' },
          { text: 'Age', value: 'age' },
        ]
      
    }
  },
  components:{
    Date,
    ContentLoader,
    Popup
  },
  
    methods:{
    getDate(e) {
      this.date = e
    },
    printdate(){
      console.log(this.date)
    },
    
    cargarClientes(){
      this.loading = true
			axios.get(process.env.VUE_APP_API_URLCLIENTS).then((res)=>{
				this.datos = res.data["datos"];
        this.loading = false
        console.log("ping")
			}).catch((error) =>{
				this.$vs.notify({
					color:'danger',
					title:'Error al cargar los clientes',
					text: error,
					iconPack: 'feather', icon:'icon-alert-circle'
				});
				console.log(error);
			});
		},
    fetchClientsbyDate(date){
      if (date==null) {
          this.errordetail = "No date was selected"
          this.dialogerror = true
          return

        }
      var date_regex = /^(19|20)\d\d[- /.](0[1-9]|1[012])[- /.](0[1-9]|[12][0-9]|3[01])$/;
      if (!(date_regex.test(date))) {
          this.dialogerror = true
          this.errordetail = "Date format not valid"
          return
      }
      this.loading = true
      console.log(date)
			axios.get(process.env.VUE_APP_API_URL_LOAD_DATA_BY_DATE,{params:{date:date}}).then((res)=>{
        this.cargarClientes()
        console.log(res.data)
			}).catch((error) =>{
				this.dialogerror = true
        this.errordetail = error
			});
		},
    loadTodayData(){
			axios.get(process.env.VUE_APP_API_URL_LOAD_TODAY_DATA).then((res)=>{
        this.cargarClientes()
        console.log(res.data)
			}).catch((error) =>{
				this.$vs.notify({
					color:'danger',
					title:'Error updating db',
					text: error,
					iconPack: 'feather', icon:'icon-alert-circle'
				});
				console.log(error);
			});
		},
    handleclick(value){
      this.currentClient = value
      this.clientName = value.name
      this.getClientTransactions(value.Cid)
      //
    },
    getClientTransactions(id){
      this.clientRecords = false
      this.clientHasNoRecords = false
      axios.get(process.env.VUE_APP_API_URLCLIENTS+id).then((res)=>{
        
        this.clientRecords = true
        this.Tdata = res.data["owner"]
        this.similarBuyer = res.data["simBuyers"]
        this.favoriteProducts = res.data["favProducts"]
			}).catch((error) =>{
        this.clientRecords = false
        this.clientHasNoRecords = true
        this.Tdata = null
        console.log(error)
				
			});

    },
    formatCurrency (value) {
    return '$' + value / 100
},
  },
  mounted(){
    this.loadTodayData()    
  }
      
}
</script>
