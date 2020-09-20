mainForm = document.getElementById("vocabList");

for (i = 0; i < 10; i++) {
    var input = document.createElement("input");
    input.type = "text"
    input.name = "input" + i
    input.id = "input" + i

    mainForm.appendChild(document.createElement('br'))
    label = document.createElement('label')
    label.htmlFor = input.name
    label.textContent = 'word' + i
    mainForm.appendChild(label)
    mainForm.appendChild(input)
    // mainForm.appendChild(document.createElement('br'))
}
