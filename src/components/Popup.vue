<template>
<div>
<v-card hover>
    <v-container>
          <v-card >
            <v-row>
            <v-col cols=6>
            <v-toolbar
              color="teal"
              dark
            >
              <v-toolbar-title>{{clientName}} Transactions</v-toolbar-title>

              <v-spacer></v-spacer>

              <v-btn icon>
                <v-icon>dots-vertical</v-icon>
              </v-btn>
            </v-toolbar>
            <v-list
            style="max-height: 250px"
            class="overflow-y-auto"
            >
      <v-list-group
        v-for="item in Tdata"
        :key="item.Tid"
        no-action
      >
            <template v-slot:activator>
              <v-list-item-content>
                <v-list-item-title >ID {{item.Tid}}</v-list-item-title>
              </v-list-item-content>
                <v-spacer />
              <v-list-item-content>
                <v-list-item-title >TOTAL {{ item.price | toCurrency }}</v-list-item-title>
              </v-list-item-content>
            </template>

            <v-list-item
              v-for="child in item.ProductIds"
              :key="child"
            >
              <v-list-item-content>
                <v-list-item-title v-text="child"></v-list-item-title>
              </v-list-item-content>
            </v-list-item>
          </v-list-group>
        </v-list>
        </v-col>
         <v-col cols=3>
           <v-toolbar
              color="teal"
              dark
            >
              <v-toolbar-title>Near Buyers</v-toolbar-title>
            </v-toolbar>
            <v-list style="max-height: 250px"
                  class="overflow-y-auto">
              <v-list-item-group
                color="primary"
              >
                <v-list-item
                  v-for="(item, i) in similarBuyer"
                  :key="i"
                >
                  <v-list-item-icon>
                    <v-icon >face</v-icon>
                  </v-list-item-icon>
                  <v-list-item-content>
                    <v-list-item-title v-text="item"></v-list-item-title>
                  </v-list-item-content>
                </v-list-item>
              </v-list-item-group>
            </v-list>
          </v-col>
          <v-col cols=3>
           <v-toolbar
              color="teal"
              dark
            >
              <v-toolbar-title>You may also like:</v-toolbar-title>
            </v-toolbar>
            <v-list style="max-height: 250px"
                  class="overflow-y-auto">
              <v-list-item-group
                color="primary"
              >
                <v-list-item
                  v-for="(item, i) in favoriteProducts"
                  :key="i"
                >
                  <v-list-item-icon>
                    <v-icon >label</v-icon>
                  </v-list-item-icon>
                  <v-list-item-content>
                    <v-list-item-title v-text="item"></v-list-item-title>
                  </v-list-item-content>
                </v-list-item>
              </v-list-item-group>
            </v-list>
          </v-col>
      </v-row>
      </v-card>
      

    </v-container>

      <v-footer dark padless>
    <v-col
      class="text-center"
      cols="12"
    >
      {{ new Date().getFullYear() }} — <strong>Vuetify</strong>
    </v-col>
  </v-footer>
      
    </v-card>
   

</div>  
</template>

<script>
//Montar el bus y recibir la info
export default {
  props:['Tdata','similarBuyer','favoriteProducts','clientName'],
    data () {
      return {
      Theader:[
        { text: 'Transaction No', value: 'Tid' },
        { text: 'Products', value: 'ProductIds' },
        { text: 'TOTAL', value: 'price' },
      ]
      }
    },
    
}
</script>