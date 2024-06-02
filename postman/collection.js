pm.test("Signup response status should be 201", function () {
pm.response.to.have.status(201);
});

pm.test("Response should contain user data", function () {
var jsonData = pm.response.json();
pm.expect(jsonData).to.have.property("username");
pm.expect(jsonData.username).to.eql("testuser");
});