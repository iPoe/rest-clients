<template>
  <div class="dashboard">
    <h1>Sin miedo al exito</h1>
    <v-container class="grey lighten-5">
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
              @click="fetchClientsbyDate(date)"
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
            <v-toolbar-title>CLIENTS</v-toolbar-title>
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
        <v-col cols=12>
          <Popup  v-bind:Tdata="Tdata" v-bind:similarBuyer="similarBuyer" />
        </v-col>
      </v-row>
    </v-container>
  </div>
</template>

<script>
import axios from 'axios'
import Date from '@/components/Date.vue'
import Popup from '../components/Popup.vue'



export default {  
  data(){
    return{
      ok:true,
      date:null,
      search: '',
      currentClient:{
            Cid: '',
            name: "",
            age: '',
        

      },
      datos:[
        
          {
            Cid: '12',
            name: "Leon",
            age: 6,
          },
      ],
      Tdata:[
        {
            "ProductIds": [
                "b0944c1f",
                "88f9107b",
                "547606cc"
            ],
            "Tid": "000060a9a4f8",
            "price": 1200,
        },
      ],
      Theader:[
        { text: 'Transaction No', value: 'Tid' },
        { text: 'Products', value: 'ProductIds' },
        { text: 'TOTAL', value: 'price' },


      ],
      similarBuyer:[
        "Sung",
        "Urion",
        "Robson",
      ],

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
			axios.get("http://localhost:9000/clients").then((res)=>{
				this.datos = res.data["datos"];
                // console.log(this.datos)
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
      const [year, month, day] = date.split('-')
      console.log(this.date)
			axios.get("http://localhost:9000/data/?date="+`${month}/${day}/${year}`).then((res)=>{
				// this.datos = res.data["datos"];
        this.cargarClientes()
        console.log(res.data)
        //Hacer el llamado a la función que obtiene la lista de clientes
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
      this.getClientTransactions(value.Cid)
      //
    },
    getClientTransactions(id){
      axios.get("http://localhost:9000/clients/"+id).then((res)=>{
        this.Tdata = res.data["owner"]
        this.similarBuyer = res.data["simBuyers"]
        //Lista de productos recomendados
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
    formatCurrency (value) {
    return '$' + value / 100
},
  },
  mounted(){
    this.cargarClientes();
  }
      
}
</script>
