<template>
    <v-row>
    <v-menu
        ref="menu"
        v-model="menu"
        :close-on-content-click="false"
        :return-value.sync="date"
        transition="scale-transition"
        offset-y
        min-width="auto"
      >
        <template v-slot:activator="{ on, attrs }">
          <v-text-field
            v-model="date"
            label="Picker in menu"
            prepend-icon="mdi-calendar"
            readonly
            v-bind="attrs"
            v-on="on"
          ></v-text-field>
        </template>
        <v-date-picker
          v-model="date"
          hint="DD/MM/YYYY format"
          no-title
          scrollable
          max = currentDate()
        >
          <v-spacer></v-spacer>
          <v-btn
            text
            color="primary"
            @click="menu = false"
          >
            Cancel
          </v-btn>
          <v-btn
            text
            color="primary"
            @click="$refs.menu.save(date);fetchClientsbyDate(date)"
          >
            OK
          </v-btn>
        </v-date-picker>
      </v-menu>
    </v-row>
</template>

<script>
import axios from 'axios'
export default {
    data: () =>  ({
      date: new Date().toISOString().substr(0, 10),
      // dateFormatted: vm.formatDate(new Date().toISOString().substr(0, 10)),
      menu: false,
      datos: null,
    }),
    methods:{
      currentDate() {
      const current = new Date();
      const date2 = `${current.getDate()}/${current.getMonth()+1}/${current.getFullYear()}`;
      console.log(date2)
      this.$emit('date', this.date)

    },
      master(){
        console.log(this.date)
        
      },
      formatDate (date) {
        if (!date) return null

        const [year, month, day] = date.split('-')
        // return `${month}/${day}/${year}`
        console.log(`${month}/${day}/${year}`)
        this.$emit('date', this.date);


      },

      fetchClientsbyDate(date){
      const [year, month, day] = date.split('-')
      this.$emit('date', date);

			axios.get("http://localhost:9000/data/?date="+`${month}/${day}/${year}`).then((res)=>{
				// this.datos = res.data["datos"];
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
    
    },
    

    mounted: function(){
      this.currentDate()
    }
  }
</script>
