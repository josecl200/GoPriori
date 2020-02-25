package main
import (
	"os"
	"bufio"
	"strings"
	"log"
	"fmt"
	"strconv"
	"sort"
)
type rule struct{
	antecedent []string
	precedent  []string
	confidence float64
}
//Prototipos para usar la funcion sort integrada en el slice de reglas
type confSort []rule
func (a confSort) Len() int           { return len(a) }
func (a confSort) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a confSort) Less(i, j int) bool { return a[i].confidence > a[j].confidence }
type lenSort []rule
func (a lenSort) Len() int           { return len(a) }
func (a lenSort) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a lenSort) Less(i, j int) bool { return (len(a[i].antecedent)+len(a[i].precedent)) < (len(a[j].antecedent)+len(a[j].precedent)) }

func Find(slice []string, val string) (int, bool) {
	for index, item := range slice {
		if item == val {
			return index,true
		}
	}
	return -1,false
}
//Devuelve los conjuntos de tamaño sizeComb dado un arreglo de itemes slice
func NCombSlices(slice []string, sizeComb uint64,sizeSlice uint64) [][]string{
	var combs [][]string
	var comb []string = make([]string,sizeComb)
	NCombSlice(&combs,slice,0,comb,sizeComb,0)
	return combs
}	
func NCombSlice(allCombs *[][]string,input []string, index uint64,output []string, sizeComb uint64,i uint64){
	if index == sizeComb{
		var out []string = make([]string,sizeComb)
		copy(out,output)
		*allCombs = append(*allCombs,out)
		return 
	}
	if i>=uint64(len(input)){
		return
	}
	output[index] = input[i]
	NCombSlice(allCombs,input,index+1,output,sizeComb,i+1)
	NCombSlice(allCombs,input,index,output,sizeComb,i+1)
}
//Verifica si los elementos de sliceHijo estan en slicePadre
func sliceInOther(slicePadre []string, sliceHijo []string) bool{
	padreMap := make(map[string]int)
	hijoMap := make(map[string]int)
	for _, padreElem := range slicePadre{
		padreMap[padreElem]++
	}
	for _, hijoElem := range sliceHijo{
		hijoMap[hijoElem]++
	}
	for hijoKey, _ := range hijoMap {
		if _,is:=padreMap[hijoKey];!is {
			return false
		}
	}
	return true
}
func SoporteConjunto (conjunto []string,transacciones [][]string,cantTrans uint64) uint64{
	var cant uint64=0
	for _,e := range transacciones{
		if sliceInOther(e,conjunto){
			cant++
		}
	}
	return cant
}
func main(){
	//El archivo con las transacciones, el soporte minimo y la confianza minima son recibidos como 1er, 2do y 3er argumentos respectivamente al llamar al script.
	file,err := os.Open(os.Args[1])
	if err != nil{
		log.Fatal(err)
	}
	defer file.Close()
	minSup,_  := strconv.ParseUint(os.Args[2],10,64)
	minConf,_ := strconv.ParseFloat(os.Args[3],64)
	var allElements []string
	var formFile []([]string)
	var transactions uint64 = 0
	scan := bufio.NewScanner(file)
	//Se sacan las transacciones al arreglo de strings formFile, tomando cada linea del archivo en formLine, y se cuenta la cantidad de transacciones
	for scan.Scan(){
		newLine := strings.Split(scan.Text()," ")
		var formLine []string
		found:=false
		for _, element := range newLine{
			found=false
			formLine = append(formLine,element)
			//Se van guardando cada elemento nuevo en un arreglo con todos los elementos existentes en una transaccion, de ya estar no se vuelven a guardar
			for _,e:=range allElements{
				if e==element{
					found=true
					break
				}
			}
			if !found{
				allElements=append(allElements,element)
			}
		}
		formFile=append(formFile,formLine)
		transactions++
	}
	//Se buscan los elementos que cumplen con el minimo soporte para formar los datasets frecuentes
	fmt.Println("Elementos que cumplen con el minimo soporte: ")
	var passed []string
	var maxSup uint64
	for i,e := range allElements{
		var strSup string
		supC:= SoporteConjunto([]string{e},formFile,transactions)
		if i==0{
			maxSup=supC
		}else if supC>maxSup{
			maxSup=supC
		}
		if supC >= minSup{
			strSup="Cumple"
			passed=append(passed,e)
		}else{
			strSup="No cumple"
		}
		fmt.Printf("Soporte para %v: %d => %s\n", e, supC,strSup)
	}
	//Se comienzan a formar datasets frecuentes a partir de los elementos que cumplen con el soporte minimo, desde un tamaño de 2 hasta que el soporte minimo no sea cumplido
	var combSize uint64 = 2
	var allPassed [][]string
	for _,e:=range passed{
		allPassed = append(allPassed,[]string{e})
	}
	sort.Strings(passed)
	//maxSup es el soporte maximo alcanzado por la ultima cantidad de itemes para los datasets
	for(maxSup>=minSup){
		newCombs:=NCombSlices(passed,combSize,uint64(len(passed)))
		passed=[]string{}
		maxSup=0
		for _,e := range newCombs{
			var strSup string
			supC:= SoporteConjunto(e,formFile,transactions)
			if supC>maxSup{
				maxSup=supC
			}
			if supC >= minSup{
				strSup="Cumple"
				allPassed=append(allPassed,e)
				for _,el:=range e{
					if _,f:=Find(passed,el);!f{
						passed=append(passed,el)
					}
				}
			}else{
				strSup="No cumple"
			}
			fmt.Printf("Soporte para %v: %d => %s\n", e, supC,strSup)
		}
		combSize++
	}
	var rules []rule
	//No es necesario, pero por precaucion se remueven itemsets malformados en los pasos anteriores (vacios y con repeticiones)
	for ind,element := range allPassed{
		if len(element)==0 || (len(element)>1 && (strings.Compare(element[0],element[1])==0)) {
			allPassed=append(allPassed[:ind], allPassed[ind+1:]...)
		}
	}
	for _,elm := range allPassed{
		if len(elm)>1 && strings.Compare(elm[0],elm[1])!=0{
			//De todos los itemsets sacados, se toman los que sean de 2 o más elementos, y se sacan subconjuntos de estos para antecedentes y procedentes
			var k uint64=1
			//se hacen subconjuntos para los antecedentes comenzando desde 1 hasta la cantidad de elementos -1 de forma que quede al menos un elemento como precedente
			for (k<=uint64(len(elm)-1)){
				antecedentes:=NCombSlices(elm,k,uint64(len(elm)))
				for _,ant:= range antecedentes{
					//se hace una copia del itemset original en cada iteracion de antecedente
					itemset:=append([]string(nil),elm...)
					for _,item:=range ant{
						//se remueven los items que forman parte del antecedente
						if indexReal,t:=Find(itemset,item);t{
							itemset = append(itemset[:indexReal],itemset[indexReal+1:]...)
						}
					}
					if len(itemset)>0{
						//Se crea una nueva regla con el antecedente y los itemes restantes del itemset
						newRule:=rule{ant,itemset,(float64(SoporteConjunto(elm,formFile,transactions))/float64(SoporteConjunto(ant,formFile,transactions)))}
						//Si la regla cumple con la minima confianza, se añade al arreglo de reglas
						if newRule.confidence>=minConf{
							rules=append(rules,newRule)
						}
						
					}
				}
				k++
			}
		}
	}
	//Se organiza el arreglo por confianza y cantidad de elementos
	sort.Sort(confSort(rules))
	sort.Sort(lenSort(rules))
	//Se imprime el arreglo de reglas, la opcion %+v nos enseñara los nombres de cada atriuto de la estructura definida anteriormente
	for _,e:=range rules{
		fmt.Printf("%+v\n",e)
	}
	
}
