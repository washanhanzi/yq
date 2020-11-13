package yqlib

import (
	"bufio"
	"container/list"
	"io"
	"os"

	yaml "gopkg.in/yaml.v3"
)

var treeNavigator = NewDataTreeNavigator(NavigationPrefs{})
var treeCreator = NewPathTreeCreator()

func readStream(filename string) (io.Reader, error) {
	if filename == "-" {
		return bufio.NewReader(os.Stdin), nil
	} else {
		return os.Open(filename) // nolint gosec
	}
}

func EvaluateStream(filename string, reader io.Reader, node *PathTreeNode, printer Printer) error {

	var currentIndex uint = 0

	decoder := yaml.NewDecoder(reader)
	for {
		var dataBucket yaml.Node
		errorReading := decoder.Decode(&dataBucket)

		if errorReading == io.EOF {
			return nil
		} else if errorReading != nil {
			return errorReading
		}
		candidateNode := &CandidateNode{
			Document: currentIndex,
			Filename: filename,
			Node:     &dataBucket,
		}
		inputList := list.New()
		inputList.PushBack(candidateNode)

		matches, errorParsing := treeNavigator.GetMatchingNodes(inputList, node)
		if errorParsing != nil {
			return errorParsing
		}
		err := printer.PrintResults(matches)
		if err != nil {
			return err
		}
		currentIndex = currentIndex + 1
	}
}

func readDocuments(reader io.Reader, filename string) (*list.List, error) {
	decoder := yaml.NewDecoder(reader)
	inputList := list.New()
	var currentIndex uint = 0

	for {
		var dataBucket yaml.Node
		errorReading := decoder.Decode(&dataBucket)

		if errorReading == io.EOF {
			switch reader := reader.(type) {
			case *os.File:
				safelyCloseFile(reader)
			}
			return inputList, nil
		} else if errorReading != nil {
			return nil, errorReading
		}
		candidateNode := &CandidateNode{
			Document: currentIndex,
			Filename: filename,
			Node:     &dataBucket,
		}

		inputList.PushBack(candidateNode)

		currentIndex = currentIndex + 1
	}
}

func EvaluateAllFileStreams(expression string, filenames []string, printer Printer) error {
	node, err := treeCreator.ParsePath(expression)
	if err != nil {
		return err
	}
	var allDocuments *list.List = list.New()
	for _, filename := range filenames {
		reader, err := readStream(filename)
		if err != nil {
			return err
		}
		fileDocuments, err := readDocuments(reader, filename)
		if err != nil {
			return err
		}
		allDocuments.PushBackList(fileDocuments)
	}
	matches, err := treeNavigator.GetMatchingNodes(allDocuments, node)
	if err != nil {
		return err
	}
	return printer.PrintResults(matches)
}

func EvaluateFileStreamsSequence(expression string, filenames []string, printer Printer) error {

	node, err := treeCreator.ParsePath(expression)
	if err != nil {
		return err
	}

	for _, filename := range filenames {
		reader, err := readStream(filename)
		if err != nil {
			return err
		}
		err = EvaluateStream(filename, reader, node, printer)
		if err != nil {
			return err
		}

		switch reader := reader.(type) {
		case *os.File:
			safelyCloseFile(reader)
		}
	}
	return nil
}

// func safelyRenameFile(from string, to string) {
// 	if renameError := os.Rename(from, to); renameError != nil {
// 		log.Debugf("Error renaming from %v to %v, attempting to copy contents", from, to)
// 		log.Debug(renameError.Error())
// 		// can't do this rename when running in docker to a file targeted in a mounted volume,
// 		// so gracefully degrade to copying the entire contents.
// 		if copyError := copyFileContents(from, to); copyError != nil {
// 			log.Errorf("Failed copying from %v to %v", from, to)
// 			log.Error(copyError.Error())
// 		} else {
// 			removeErr := os.Remove(from)
// 			if removeErr != nil {
// 				log.Errorf("failed removing original file: %s", from)
// 			}
// 		}
// 	}
// }

// // thanks https://stackoverflow.com/questions/21060945/simple-way-to-copy-a-file-in-golang
// func copyFileContents(src, dst string) (err error) {
// 	in, err := os.Open(src) // nolint gosec
// 	if err != nil {
// 		return err
// 	}
// 	defer safelyCloseFile(in)
// 	out, err := os.Create(dst)
// 	if err != nil {
// 		return err
// 	}
// 	defer safelyCloseFile(out)
// 	if _, err = io.Copy(out, in); err != nil {
// 		return err
// 	}
// 	return out.Sync()
// }

func safelyFlush(writer *bufio.Writer) {
	if err := writer.Flush(); err != nil {
		log.Error("Error flushing writer!")
		log.Error(err.Error())
	}

}
func safelyCloseFile(file *os.File) {
	err := file.Close()
	if err != nil {
		log.Error("Error closing file!")
		log.Error(err.Error())
	}
}
