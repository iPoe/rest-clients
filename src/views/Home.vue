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
          <v-data-table
            :headers="headers"
            :items="datos"
            :items-per-page="5"
            :sort-by="['name', 'age']"
            class="elevation-1"
            @click:row="handleclick"
            
          >
           
         
          </v-data-table>

        </v-col>
        
        <v-col cols="7">
          <v-card >
            <v-card-text>
         
          <v-row >
            <v-col cols=4 >
          <v-text-field
            v-model ="currentClient.Cid"
            label="ID"
            outlined
            readonly
          ></v-text-field>
          
        </v-col>
        <v-col cols=4>
          <v-text-field
            v-model ="currentClient.name" 
            label="Name"
            outlined
            readonly
          ></v-text-field>
        </v-col>
        <v-col cols=2>
          <v-text-field
            v-model ="currentClient.age"
            label="Age"
            outlined
            readonly
          ></v-text-field>
        </v-col>

          </v-row>
          <v-row>
            <v-col cols=12>
              <v-toolbar
            class="mb-2"
            color="purple"
            dark
            flat
          >
            <v-toolbar-title> TRANSACTIONS</v-toolbar-title>
          </v-toolbar>
              <v-data-table
              dense
            :headers="Theader"
            :items="Tdata"
            :items-per-page="3"
            class="elevation-1"

          >
          <template v-slot:[`item.price`]="{ item }">
            <span>{{ item.price | toCurrency }}</span>
          </template>
          </v-data-table>
            </v-col>
            
          </v-row>
          <v-row>
            <v-sheet
    class="mx-auto"
    max-width="500"
  >
    <v-slide-group
      multiple
      show-arrows
    >
      <v-slide-item
        v-for="n in similarBuyer"
        :key="n"
        v-slot="{ active, toggle }"
      >
        <v-btn
          class="mx-2"
          :input-value="active"
          active-class="purple white--text"
          depressed
          rounded
          @click="toggle"
        >
           {{ n }}
        </v-btn>
      </v-slide-item>
    </v-slide-group>
  </v-sheet>
          </v-row>
            </v-card-text>
          </v-card>

        </v-col>
      </v-row>
      

    </v-container>
  </div>
</template>

<script>
import axios from 'axios'
import Date from '@/components/Date.vue'


export default {  
  data(){
    return{
      ok:true,
      date:null,
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
        { text: 'Tid', value: 'Tid' },
        { text: 'Products', value: 'ProductIds' },
        { text: 'Price', value: 'price' },


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
    Date
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
    formatCurrency (value) {
    return '$' + value / 100
},
  },
  mounted(){
    this.cargarClientes();
  }
      
}
</script>
