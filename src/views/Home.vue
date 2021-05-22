<template>
  <div class="dashboard">
    <h1>Sin miedo al exito</h1>
    <v-container class="grey lighten-5">
      <v-row justify="space-between">
        <v-col md ="2" align='center'>
          <Date v-on:date="getDate($event)" />
       <v-btn
        color="primary"
        elevation="2"
        @click="fetchClientsbyDate(date)"
      >Search
      <v-icon right>search</v-icon>
      </v-btn>

        </v-col>
      </v-row>
      <v-row justify="center">
        <v-col md ="6" >
        <v-toolbar
            class="mb-2"
            color="indigo darken-5"
            dark
            flat
          >
            <v-toolbar-title>Users List</v-toolbar-title>
          </v-toolbar>
          <v-data-table
            :headers="headers"
            :items="datos"
            :items-per-page="5"
            :sort-by="['name', 'age']"
            class="elevation-1"
            @click:row="handleclick"
          ></v-data-table>

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
      date:null,
      datos:[
        
          {
            Cid: '12',
            name: "Leon",
            age: 6,
          },
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
      console.log(value)
      this.$router.push('/about')

    }
  },
  mounted(){
    this.cargarClientes();
  }
      
}
</script>
