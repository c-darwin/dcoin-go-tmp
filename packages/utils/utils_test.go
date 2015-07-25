package utils

import (
	"testing"
	"fmt"
)

func TestDecryptCFB(t *testing.T) {
	decr, err := DecryptCFB([]byte("CPd2fWi3x4JhN8Ar"), HexToBin([]byte("17b803ff49ff86496fe85df3430ead37455f4d14273e58dc6db08874229bd69e21a9ba611cb1ee3d668656a82bacfb8f9ab3de4557a6612cd45e6011063fca89c188d4b5b3e5328a3f18ce6bef3c537b3e6e31c299ccdd2ba730570f2befb5d34d9377bba92b9e9b29914fe93093f2fc98a28fd465ddaa5d78f980a993536365cb097d0e1067a8905abe66bb5c5462a4715a1940515d671f04eb3a0ab37ab4e6ae6730468f9b98ff9fb7047c8b22507121908fef0b333d1c199c2c48c142976ab8ec72364af13e8625adba53c516dba867d36e7f2a256b72fcdbc08e6c35bc0ccadc0828e8f6189f5c31519e4e0010baf530b5f5f4913396838a2f732a69c5142c2b5ae2db88e85abfd17165490f41165fe7807f0d8e496f4074659f83b4f6ca2e8b45e85e40d2765442d00ba40d5f61e00cf4c44ba8111525165b07e2eaa448ab1c9aeb43ad761bf9c33c740a1b62d628ac6cde39970dfcf379f2c0fbc95f0e5825ae9a965b53d34a374a6a6c7aab74b8ff6c8c99483ac667875d33a5e3e7d5eaef340c19cbaab11d1a39891d27a1783067e1f7e77d33c3569446ec657c52814adcac5e422562d6c77f026cd11215ce8f5dbc0462bc2f4f427f25f4bbf7ee5e0d23f0bf6091cde0626ced4236a375f43affe2ee57228461a4fb5bb331a8cf5c8e994a1917730a23f4b932636b0bcb22c640f9b9b77ff6a3d844ea3f4b827a0a0437165bd3ff51fb2c647b79fa51d913a31ac687155f23fd8e4ef22af40f6c83fb98a8a734435b07563343a66c4ae4af9bbe")), []byte("244f3aeac7701f78818401d093a82fad"))
	fmt.Printf("decr %x\n", decr)
	if err!=nil {
		t.Error(err)
	}
}


func TestFormatQueryArgs(t *testing.T) {
  	newQuery, newArgs := FormatQueryArgs("SELECT id, to_user_id FROM credits WHERE from_user_id = ? AND currency_id = ? AND tx_block_id = ? AND tx_hash = [hex] AND del_block_id = 0 ORDER BY time DESC", "sqlite", []interface {}{1, 1, 1, "FFFFFFFFFFFF"}...)
	fmt.Println(newQuery)
	fmt.Println(newArgs)

}