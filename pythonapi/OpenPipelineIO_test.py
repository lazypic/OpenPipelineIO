from OpenPipelineIO import opio

print(opio.Project("mkk3"))
print(opio.Shot("mkk3", "JYW_0200"))
print(opio.Projects())
print(opio.Seqs("TEMP"))
print(opio.Shots("TEMP", "SS"))
print(opio.searchwordItems("TEMP", "fx"))
