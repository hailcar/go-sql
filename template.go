package main

import (
	"fmt"
)

func GenerateTemplate(input string) {
	capitalizedInput, uncapitalizedInput := getCasesOfString(input)

	fmt.Printf(
		`
use app\common\utils\ObjectUtils; 


class %[1]sDao
{
    public function findById(int $id):?%[1]s{
        $model = new %[1]sModel();
        $result=$model->find($id);
        return ObjectUtils::getObjectFromModel(%[1]s::class,$result);
    }

    public function delete(int $id):bool{
        if($this->findById($id)==null){
            return true;
        }else{
            $bool=%[1]sModel::find($id)->delete();
            return $bool;
        }
    }
    public function save(%[1]s $input):?%[1]s{
        $model = new %[1]sModel();
        if(isset($input->id)){
            $model->find($input->id)->save((array)$input);
            $model=$model->find($input->id);
        }else{
            $model->save((array)$input);
        }
        return ObjectUtils::getObjectFromModel(%[1]s::class, $model);
    }
} 


class %[1]sManager
{
    private %[1]sDao $%[2]sDao;

    public function __construct(
       %[1]sDao $%[2]sDao
    )
    {
        $this->%[2]sDao = $%[2]sDao;
    }

    public function save(%[1]s $%[2]s):?%[1]s{
    return $this->%[2]sDao->save($%[2]s);
    }
    public function findById(int $id):?%[1]s{
        return $this->%[2]sDao->findById($id);
    }
    public function delete(int $id):bool{
        return $this->%[2]sDao->delete($id);
    }
}
`, capitalizedInput, uncapitalizedInput)
}
