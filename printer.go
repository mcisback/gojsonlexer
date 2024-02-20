package main

import "fmt"

func printRepeatStr(str string, level int) {
	// fmt.Print(level)

	for i := 0; i < level; i++ {
		fmt.Print(str)
	}
}

func printParser(node *JsonNode, level int) {
	printRepeatStr("-", level)

	// TODO Print Pairs

	if node.Type == JSON_NODE_ARRAY {
		fmt.Println("Array ", "(", len(node.Children), ") ", "([")
	} else if node.Type == JSON_NODE_OBJECT {
		fmt.Println("Object ", "(", len(node.Children), ") ", "{")
	} else {
		fmt.Println("Node(t:", node.Type, ", v:", node.Value, ", c:", len(node.Children), ")")
	}

	if len(node.Children) > 0 {
		if node.Type == JSON_NODE_ARRAY {
			for _, child := range node.Children {
				printParser(&child, level+1)
			}
		}
	}

	if node.Type == JSON_NODE_OBJECT {
		object := node.Value.(JsonObject)

		for key := range object.Pairs {
			pair := object.Pairs[key]

			printParser(&pair, level+1)
		}
	}

	if node.Type == JSON_NODE_ARRAY {
		printRepeatStr("-", level)

		fmt.Println("])")
	} else if node.Type == JSON_NODE_OBJECT {
		printRepeatStr("-", level)

		fmt.Println("}")
	}
}

func printJsonValue(node *JsonNode, isAtomic bool) {
	value := node.Value.(JsonValue)

	switch value.Type {
	case JSON_VALUE_BOOLEAN, JSON_VALUE_NUMBER, JSON_VALUE_NULL:
		fmt.Print(value.Value)
	case JSON_VALUE_STRING:
		if !isAtomic {
			fmt.Print("\"")
		}

		fmt.Print(value.Value)

		if !isAtomic {
			fmt.Print("\"")
		}
	}

	if !isAtomic {
		fmt.Println(",")
	}
}

func printJson(node *JsonNode, level int, key string) {
	separator := "\t"

	printRepeatStr(separator, level)

	if key != "" {
		fmt.Print("\"", key, "\": ")
	}

	// TODO Print Pairs

	if node.Type == JSON_NODE_ARRAY {
		fmt.Println("[")
	} else if node.Type == JSON_NODE_OBJECT {
		fmt.Println("{")
	} else if node.Type == JSON_NODE_VALUE {
		printJsonValue(node, false)
	}

	if len(node.Children) > 0 {
		if node.Type == JSON_NODE_ARRAY {
			for _, child := range node.Children {
				printJson(&child, level+1, "")
			}
		}
	}

	if node.Type == JSON_NODE_OBJECT {
		object := node.Value.(JsonObject)

		for key := range object.Pairs {
			pair := object.Pairs[key]

			printJson(&pair, level+1, key)
		}
	}

	if node.Type == JSON_NODE_ARRAY {
		printRepeatStr(separator, level)

		fmt.Println("],")
	} else if node.Type == JSON_NODE_OBJECT {
		printRepeatStr(separator, level)

		fmt.Println("},")
	}
}
